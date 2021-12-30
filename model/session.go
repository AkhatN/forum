package model

//Session struct
type Session struct {
	UUID   string
	AuthID int
}

//WriteUUIDtoDataBase writes UUID of a user into session table in database
func (s *Session) WriteUUIDtoDataBase() error {
	_, err := Db.Exec("INSERT INTO session (uuid, auth_id) VALUES (?, ?)", s.UUID, s.AuthID)
	if err != nil {
		return err
	}

	return nil
}

//DeleteSessionFromDB deletes session from db
func (s *Session) DeleteSessionFromDB() error {
	_, err := Db.Exec("DELETE FROM session WHERE auth_id = ?", s.AuthID)
	if err != nil {
		return err
	}
	return nil
}

//ReadAuthIdbyUUID reads user_id by uuid provided in a request
func (s *Session) ReadAuthIdbyUUID() error {
	row := Db.QueryRow("SELECT auth_id FROM session WHERE uuid = ?", s.UUID)
	err := row.Scan(&s.AuthID)
	if err != nil {
		return err
	}
	return nil
}
