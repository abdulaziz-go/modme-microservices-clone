package server

import (
	"education-service/config"
	"education-service/internal/clients"
	"education-service/internal/repository"
	"education-service/internal/service"
	"education-service/internal/utils"
	"education-service/migrations"
	"education-service/proto/pb"
	"fmt"
	_ "github.com/lib/pq"
	"github.com/robfig/cron/v3"
	"google.golang.org/grpc"
	"log"
	"net"
	"strconv"
	"time"
)

func RunServer() {
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}
	db, err := repository.NewPostgresDB(&cfg.Database)
	if err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}
	defer db.Close()
	migrations.SetUpMigrating(cfg.Database.Action, db)

	userClient, err := clients.NewUserClient(cfg.Grpc.UserService.Address)
	if err != nil {
		log.Fatalf("error %v", err)
	}

	financeClientChanForAttendance := make(chan *clients.FinanceClient)
	financeClientChanForStudent := make(chan *clients.FinanceClient)

	go func() {
		time.Sleep(2 * time.Second)
		var client *clients.FinanceClient
		for {
			fmt.Println(cfg.Grpc.FinanceService.Address)
			client, err = clients.NewFinanceClient(cfg.Grpc.FinanceService.Address)
			if err == nil {
				log.Println("Connected to Finance Service successfully.")
				financeClientChanForAttendance <- client
				close(financeClientChanForAttendance)
				break
			}
			log.Printf("Waiting for Finance Service...")
			time.Sleep(2 * time.Second)
		}
	}()

	go func() {
		time.Sleep(2 * time.Second)
		var client *clients.FinanceClient
		for {
			fmt.Println(cfg.Grpc.FinanceService.Address)
			client, err = clients.NewFinanceClient(cfg.Grpc.FinanceService.Address)
			if err == nil {
				log.Println("Connected to Finance Service successfully.")
				financeClientChanForStudent <- client
				close(financeClientChanForStudent)
				break
			}
			log.Printf("Waiting for Finance Service...")
			time.Sleep(2 * time.Second)
		}
	}()

	roomRepo := repository.NewRoomRepository(db)
	roomService := service.NewRoomService(roomRepo)
	courseRepo := repository.NewCourseRepository(db)
	courseService := service.NewCourseService(courseRepo)
	groupRepo := repository.NewGroupRepository(db, userClient)
	groupService := service.NewGroupService(groupRepo)
	attendanceRepo := repository.NewAttendanceRepository(db, financeClientChanForAttendance)
	attendanceService := service.NewAttendanceService(attendanceRepo)
	studentRepo := repository.NewStudentRepository(db, userClient, financeClientChanForStudent)
	studentService := service.NewStudentService(studentRepo)
	companyRepo := repository.NewCompanyRepository(db, userClient)
	companyService := service.NewCompanyService(companyRepo)
	tarrifRepo := repository.NewTariffRepository(db)
	tarrifService := service.NewTariffService(tarrifRepo)
	companyFinanceRepo := repository.NewCompanyFinanceRepository(db)
	companyFinanceService := service.NewCompanyFinanceService(companyFinanceRepo)
	smsRepository := repository.NewSmsRepository(db)
	smsService := service.NewSmsService(smsRepository)
	lis, err := net.Listen("tcp", ":"+strconv.Itoa(cfg.Server.Port))
	if err != nil {
		log.Fatalf("Failed to listen on port %v: %v", cfg.Server.Port, err)
	}

	grpcServer := grpc.NewServer(
		grpc.UnaryInterceptor(utils.RecoveryInterceptor),
	)
	pb.RegisterRoomServiceServer(grpcServer, roomService)
	pb.RegisterCourseServiceServer(grpcServer, courseService)
	pb.RegisterGroupServiceServer(grpcServer, groupService)
	pb.RegisterAttendanceServiceServer(grpcServer, attendanceService)
	pb.RegisterStudentServiceServer(grpcServer, studentService)
	pb.RegisterCompanyServiceServer(grpcServer, companyService)
	pb.RegisterTariffServiceServer(grpcServer, tarrifService)
	pb.RegisterCompanyFinanceServiceServer(grpcServer, companyFinanceService)
	pb.RegisterSmsServiceServer(grpcServer, smsService)

	c := cron.New()
	_, err = c.AddFunc("10 1 1 * *", func() {
		fmt.Println("Running student balance taker ....")
		studentRepo.StudentBalanceTaker()
		fmt.Println("Completed student balance taker")
	})

	if err != nil {
		log.Fatalf("Failed to schedule cron job: %v", err)
	}
	c.Start()

	go func() {
		select {}
	}()

	log.Printf("Server listening on port %v", cfg.Server.Port)
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}
