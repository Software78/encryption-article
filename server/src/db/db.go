
package db

import "gorm.io/gorm"

type Database interface { // Define the interface
    First(dest interface{}, conds ...interface{}) *gorm.DB
    Create(value interface{}) *gorm.DB
    Save(value interface{}) *gorm.DB
    Delete(value interface{}, conds ...interface{}) *gorm.DB
    Find(dest interface{}, conds ...interface{}) *gorm.DB
	Where(query interface{}, args ...interface{}) *gorm.DB
	AutoMigrate(dst ...interface{}) error
    // Add all the other methods you use from gorm.DB
}

type GormDB struct {
    *gorm.DB
}

func NewGormDB(db *gorm.DB) Database {
    return &GormDB{db}
}

func (g *GormDB) First(dest interface{}, conds ...interface{}) *gorm.DB {
    return g.DB.First(dest, conds...)
}

func (g *GormDB) Create(value interface{}) *gorm.DB {
    return g.DB.Create(value)
}

func (g *GormDB) Save(value interface{}) *gorm.DB {
    return g.DB.Save(value)
}

func (g *GormDB) Delete(value interface{}, conds ...interface{}) *gorm.DB {
    return g.DB.Delete(value, conds...)
}

func (g *GormDB) Find(dest interface{}, conds ...interface{}) *gorm.DB {
    return g.DB.Find(dest, conds...)
}

func (g *GormDB) AutoMigrate(dst ...interface{}) error {
    return g.DB.AutoMigrate(dst...)
}

func (g *GormDB) Where(query interface{}, args ...interface{}) *gorm.DB {
	return g.DB.Where(query, args...)
}
