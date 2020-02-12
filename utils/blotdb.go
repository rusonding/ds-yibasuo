package utils

import (
  "bytes"
  "fmt"
  "github.com/boltdb/bolt"
  "log"
)

var (
  //DbUrl = "C:\\iceblue\\go\\src\\ds-yibasuo\\web\\ds-yibasuo.db"
  DbUrl = "ds-yibasuo.db"
  Db    *BoltManager
)

func BlotInit() {
  Db, _ = NewBoltManager(DbUrl, []string{"host"})
}

//BlotDB的管理类
type BoltManager struct {
  db *bolt.DB
}

//创建库管理,并生成Bucket
func NewBoltManager(dbPath string, bucket []string) (*BoltManager, error) {
  //  bolt.Open("my.db", 0600, &bolt.Options{Timeout: 1 * time.Second,ReadOnly: true})
  db, err := bolt.Open(dbPath, 0644, nil)
  if err != nil {
    return nil, err
  }
  err = db.Update(func(tx *bolt.Tx) error {
    for _, v := range bucket {
      _, err := tx.CreateBucketIfNotExists([]byte(v))
      if err != nil {
        return err
      }
    }
    return nil
  })

  if err != nil {
    return nil, err
  }
  return &BoltManager{db}, nil
}

//关库数据库
func (m *BoltManager) Close() error {
  return m.db.Close()
}

//移除Bucket
func (m *BoltManager) RemoveBucket(bucketName string) (err error) {
  err = m.db.Update(func(tx *bolt.Tx) error {
    return tx.DeleteBucket([]byte(bucketName))
  })
  return err
}

//组Bucket增加值
func (m *BoltManager) Add(bucketName string, val []byte) (id uint64, err error) {
  err = m.db.Update(func(tx *bolt.Tx) error {
    b := tx.Bucket([]byte(bucketName))
    id, _ = b.NextSequence() //sequence uint64
    bBuf := fmt.Sprintf("%d", id)
    return b.Put([]byte(bBuf), val)
  })
  return
}

//遍历
func (m *BoltManager) Select(bucketName string) (res map[string]string, err error) {
  res = make(map[string]string)
  _ = m.db.View(func(tx *bolt.Tx) error {
    b := tx.Bucket([]byte(bucketName))
    _ = b.ForEach(func(k, v []byte) error {
      //log.Printf("key=%s, value=%s\n", string(k), v)
      res[string(k)] = string(v)
      return nil
    })
    return nil
  })
  return res, nil
}

//移除指定Bucket中指定ID
func (m *BoltManager) RemoveID(bucketName string, id []byte) error {
  err := m.db.Update(func(tx *bolt.Tx) error {
    b := tx.Bucket([]byte(bucketName))
    return b.Delete(id)
  })
  return err
}

//移除指定Bucket中指定Val
func (m *BoltManager) RemoveVal(bucketName string, val []byte) (err error) {
  var arrID []string
  arrID = make([]string, 1)
  err = m.db.View(func(tx *bolt.Tx) error {
    b := tx.Bucket([]byte(bucketName))
    c := b.Cursor()
    for k, v := c.First(); k != nil; k, v = c.Next() {
      log.Printf("key=%s, value=%s\n", k, string(v))
      if bytes.Compare(v, val) == 0 {
        arrID = append(arrID, string(k))
      }
    }
    return nil
  })

  err = m.db.Update(func(tx *bolt.Tx) error {
    b := tx.Bucket([]byte(bucketName))
    for _, v := range arrID {
      _ = b.Delete([]byte(v))
      log.Println("Del k:", v)
    }
    return nil
  })

  return err
}

//查找指定值
//func (m *BoltManager) SelectVal(bucketName string, val []byte) (arr []string, err error) {
//  arr = make([]string, 0, 1)
//  err = m.db.View(func(tx *bolt.Tx) error {
//    c := tx.Bucket([]byte(bucketName)).Cursor()
//    for k, v := c.First(); k != nil; k, v = c.Next() {
//      if bytes.Compare(v, val) == 0 {
//        arr = append(arr, string(k))
//      }
//    }
//    return nil
//  })
//  return arr, err
//}

func (m *BoltManager) SelectVal(bucketName string, val []byte) (arr []string, err error) {
  arr = make([]string, 0, 1)
  err = m.db.View(func(tx *bolt.Tx) error {
    b := tx.Bucket([]byte(bucketName))
    get := b.Get(val)
    if string(get) != "" {
      arr = append(arr, string(get))
    }
    return nil
  })
  return arr, err
}

//在事务中，移除指定Bucket中指定Val
func (m *BoltManager) RemoveValTransaction(bucketName string, val []byte) (count int, err error) {
  arrID, err1 := m.SelectVal(bucketName, val)
  if err1 != nil {
    return 0, err1
  }
  count = len(arrID)
  if count == 0 {
    return count, nil
  }

  tx, err1 := m.db.Begin(true)
  if err1 != nil {
    return count, err1
  }
  b := tx.Bucket([]byte(bucketName))
  for _, v := range arrID {
    if err = b.Delete([]byte(v)); err != nil {
      log.Printf("删除ID(%s)失败! 执行回滚. err:%s \r\n", v, err)
      _ = tx.Rollback()
      return
    }
    log.Println("删除ID(", v, ")成功!")
  }
  err = tx.Commit()
  return
}
