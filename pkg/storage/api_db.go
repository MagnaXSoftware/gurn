//go:build !nodb

package storage

import "fmt"

func (u Urn) FindByName(name string) *Urn {
	db := GetDatabase()

	n := Urn{}
	tx := db.Model(&u).First(&n, "name = ?", name)
	if tx.Error != nil {
		return nil
	}
	return &n
}

func (u *Urn) Save() error {
	db := GetDatabase()

	tx := db.Save(u)
	return tx.Error
}

func (u *Urn) Delete() error {
	if u.ID == 0 {
		return fmt.Errorf("urn has no ID, can't be removed from storage")
	}

	db := GetDatabase()
	tx := db.Delete(u)
	return tx.Error
}

func (a Access) FindByUrn(urn string) []Access {
	db := GetDatabase()

	var accesses []Access
	tx := db.Where("urn = ?", urn).Find(&accesses)
	if tx.Error != nil {
		return []Access{}
	}
	return accesses

}

func (a *Access) Save() error {
	db := GetDatabase()

	tx := db.Create(a)
	return tx.Error
}
