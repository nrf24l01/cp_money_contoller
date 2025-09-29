package schemas

type TaskType struct {
	TaskName     string          `json:"task_name" bson:"task_name" validate:"required,min=1,max=100"`
	TaskType     string          `json:"task_type" bson:"task_type" validate:"required,min=1,max=50" gorm:"uniqueIndex"`
	TaskTemplate []TaskTypeField `json:"task_template" bson:"task_template" validate:"required,dive"`
}

type TaskTypeField struct {
	FieldName        string `json:"field_name" bson:"field_name" validate:"required,min=1,max=50"`
	FieldHint        string `json:"field_hint" bson:"field_hint" validate:"required,max=200"`
	FieldDescription string `json:"field_description" bson:"field_description" validate:"max=500"`
}

type TaskTypeWithUUID struct {
	UUID string `json:"uuid" bson:"uuid" validate:"required,uuid"`
	TaskType
}

type CreateTaskTypeRequest struct {
	TaskName     string          `json:"task_name" bson:"task_name" validate:"required,min=1,max=100"`
	TaskType     string          `json:"task_type" bson:"task_type" validate:"required,min=1,max=50"`
	TaskTemplate []TaskTypeField `json:"task_template" bson:"task_template" validate:"required,dive"`
}

type CreateTaskTypeResponse struct {
	UUID         string          `json:"uuid" bson:"uuid" validate:"required,uuid"`
	TaskName     string          `json:"task_name" bson:"task_name" validate:"required,min=1,max=100"`
	TaskType     string          `json:"task_type" bson:"task_type" validate:"required,min=1,max=50"`
	TaskTemplate []TaskTypeField `json:"task_template" bson:"task_template" validate:"required,dive"`
}