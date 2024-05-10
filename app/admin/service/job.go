package service

import (
	"errors"
	"github.com/MarchGe/go-admin-server/app/admin/model"
	"github.com/MarchGe/go-admin-server/app/admin/service/dto/req"
	"github.com/MarchGe/go-admin-server/app/admin/service/dto/res"
	"github.com/MarchGe/go-admin-server/app/common/E"
	"github.com/MarchGe/go-admin-server/app/common/database"
	"gorm.io/gorm"
	"time"
)

var _jobService = &JobService{}

type JobService struct {
}

func GetJobService() *JobService {
	return _jobService
}

func (s *JobService) CreateJob(info *req.JobUpsertReq) error {
	job := s.toModel(info)
	job.UpdateTime = time.Now()
	job.CreateTime = time.Now()
	err := database.GetMysql().Transaction(func(tx *gorm.DB) error {
		return tx.Save(job).Error
	})
	return err
}

func (s *JobService) toModel(info *req.JobUpsertReq) *model.Job {
	return &model.Job{
		Name:        info.Name,
		SortNum:     info.SortNum,
		Description: info.Description,
	}
}

func (s *JobService) UpdateJob(id int64, info *req.JobUpsertReq) error {
	job, _ := s.FindOneById(id)
	if job == nil {
		return E.Message("操作的岗位不存在")
	}
	s.copyProperties(info, job)
	job.UpdateTime = time.Now()
	err := database.GetMysql().Transaction(func(tx *gorm.DB) error {
		return tx.Save(job).Error
	})
	return err
}

func (s *JobService) FindOneById(id int64) (*model.Job, error) {
	m := &model.Job{}
	err := database.GetMysql().First(m, id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, E.Message("查询的记录不存在")
		} else {
			return nil, err
		}
	}
	return m, nil
}

func (s *JobService) copyProperties(info *req.JobUpsertReq, job *model.Job) {
	job.Name = info.Name
	job.SortNum = info.SortNum
	job.Description = info.Description
}

func (s *JobService) DeleteJob(id int64) error {
	err := database.GetMysql().Transaction(func(tx *gorm.DB) error {
		if err := tx.Delete(&model.Job{}, id).Error; err != nil {
			return err
		}
		if err := tx.Where("job_id = ?", id).Delete(&model.UserJob{}).Error; err != nil {
			return err
		}
		return nil
	})
	return err
}

func (s *JobService) PageList(keyword string, page, pageSize int) (*res.PageableData[*model.Job], error) {
	jobs := make([]*model.Job, 0)
	pageableData := &res.PageableData[*model.Job]{}
	db := database.GetMysql().Model(&model.Job{})
	if keyword != "" {
		db.Where("name like ?", "%"+keyword+"%")
	}
	var count int64
	err := db.Count(&count).Order("sort_num").Offset(pageSize * (page - 1)).Limit(pageSize).Find(&jobs).Error
	if err != nil {
		return nil, err
	}
	pageableData.List = jobs
	pageableData.Total = count
	return pageableData, nil
}
