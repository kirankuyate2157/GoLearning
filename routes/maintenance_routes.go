package routes

import (
	controllers "app/controllers/maintenance"

	"github.com/labstack/echo/v4"
)

func MaintenanceRoutes(e *echo.Echo) {
	// Base group: /api/Maintenances

	g := e.Group("/api/maintenance")

	// Maintenance Type Endpoints
	g.GET("/type", controllers.ListMaintenanceTypes)
	g.POST("/type", controllers.CreateMaintenanceType)
	g.GET("/type/:uid", controllers.GetMaintenanceType)
	g.PATCH("/type/:uid", controllers.UpdateMaintenanceType)
	g.DELETE("/type/:uid", controllers.DeleteMaintenanceType)

	// Maintenance Template Endpoints
	g.GET("/templates", controllers.ListMaintenanceTemplate)
	g.POST("/templates", controllers.CreateMaintenanceTemplate)
	g.GET("/templates/:uid", controllers.GetMaintenanceTemplate)
	g.PATCH("/templates/:uid", controllers.UpdateMaintenanceTemplate)
	g.DELETE("/templates/:uid", controllers.DeleteMaintenanceTemplate)

	// Maintenance Template Machine Endpoints
	g.GET("/templates/:template_uid/machines", controllers.GetMaintenanceTemplateMachineList)
	g.POST("/templates/:template_uid/machines", controllers.CreateMaintenanceTemplateMachine)
	g.DELETE("/templates/:template_uid/machines/:machine_uid", controllers.DeleteMaintenanceTemplateMachine)

	// Maintenance Template Procedure Endpoints
	g.GET("/templates/:uid/procedures", controllers.ListTemplateProcedures)
	// g.GET("/templates//procedures", controllers.ListProceduresForATemplate)
	g.POST("/templates/:uid/procedures", controllers.CreateTemplateProcedure)
	g.GET("/templates/procedures/:uid", controllers.GetTemplateProcedure)
	g.PATCH("/templates/procedures/:uid", controllers.UpdateTemplateProcedure)
	g.DELETE("/templates/procedures/:uid", controllers.DeleteTemplateProcedure)

	// Maintenance Procedure Steps Endpoints
	g.GET("/procedures/:procedure_uid/steps", controllers.ListProcedureSteps)
	g.POST("/procedures/:procedure_uid/steps", controllers.CreateProcedureSteps)
	g.DELETE("/procedures/steps/:uid", controllers.DeleteProcedureSteps)
	g.PATCH("/procedures/steps/:uid", controllers.UpdateProcedureStep)

	// Maintenance Procedure Step Question Endpoints
	g.GET("/procedures/steps/:uid/questions", controllers.ListProcedureStepQuestions)
	g.POST("/procedures/steps/:uid/questions", controllers.CreateProcedureStepQuestion)
	g.PATCH("/procedures/steps/questions/:uid", controllers.UpdateProcedureStepQuestion)
	g.DELETE("/procedures/steps/questions/:uid", controllers.DeleteProcedureStepQuestion)

	// Maintenance Procedure Step Question Answer Endpoints
	g.GET("/procedures/steps/questions/:uid/answers", controllers.ListQuestionAnswers)
	g.POST("/procedures/steps/questions/:uid/answers", controllers.CreateQuestionAnswer)
	g.PATCH("/procedures/steps/questions/answers/:uid", controllers.UpdateQuestionAnswer)
	g.DELETE("/procedures/steps/questions/answers/:uid", controllers.DeleteQuestionAnswer)

	// Maintenance Schedule Endpoints
	g.GET("/schedules", controllers.ListMaintenanceSchedules)
	g.POST("/schedules", controllers.CreateMaintenanceSchedule)
	// g.GET("/schedules/:uid", controllers.GetMaintenanceSchedule)
	// g.PATCH("/schedules/:uid", controllers.UpdateMaintenanceSchedule)
	g.DELETE("/schedules/:uid", controllers.DeleteMaintenanceSchedule)

	// Maintenance Schedule Question Endpoints
	g.GET("/schedules/:maintenance_schedule_uid/observations", controllers.ListMaintenanceData)
	g.POST("/schedules/:maintenance_schedule_uid/observations", controllers.CreateMaintenanceScheduleObservation)
	g.DELETE("/schedules/:maintenance_schedule_uid/observations/:procedure_step_question_uid", controllers.DeleteMaintenanceScheduleQuestion)

}
