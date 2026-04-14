// Copyright 2026 Jakkaphat Chalermphanaphan

// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at

//     https://www.apache.org/licenses/LICENSE-2.0

// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package store

import (
	"errors"
	"time"

	"github.com/gxjakkap/reception/models"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type PendingStore struct {
	db *gorm.DB
}

func NewPendingStore(db *gorm.DB) *PendingStore {
	return &PendingStore{
		db: db,
	}
}

func (s *PendingStore) Create(m *models.Pending) (error, int8) {
	var p models.Pending
	res := s.db.Where("user_id = ?", m.UserID).Limit(1).Find(&p)
	if res.Error != nil {
		return res.Error, -1
	}

	if res.RowsAffected > 0 && p.ExpiredAt.Before(time.Now()) {
		return errors.New("pending request already exists"), 1
	}

	if err := s.db.Create(m).Error; err != nil {
		return err, 2
	}
	return nil, 0
}

func (s *PendingStore) Fulfill(i uint) (models.Pending, error) {
	var pen models.Pending
	err := s.db.Clauses(clause.Returning{}).Where("id = ?", i).Delete(&pen).Error
	if err != nil {
		return models.Pending{}, err
	}
	return pen, nil
}

func (s *PendingStore) DeleteExpired() error {
	if err := s.db.Where("expired_at < ?", time.Now()).Delete(&models.Pending{}).Error; err != nil {
		return err
	}
	return nil
}

func (s *PendingStore) UserInPendingList(u string) (bool, error) {
	var p models.Pending
	res := s.db.Where("user_id = ?", u).Limit(1).Find(&p)
	if res.Error != nil {
		return false, res.Error
	}
	return res.RowsAffected > 0, nil
}

func (s *PendingStore) GetPendingByUser(u string) (models.Pending, error) {
	var pen models.Pending
	err := s.db.Where("user_id = ?", u).Limit(1).Find(&pen).Error
	if err != nil {
		return models.Pending{}, err
	}
	return pen, nil
}
