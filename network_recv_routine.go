func (d *down) recv() {

  defer d.conn.Close()

	for {
		buf := make([]byte, 2048, 2048)
		d.conn.SetReadDeadline(time.Now().Add(1 * time.Second))
		n, _ := d.conn.Read(buf)

		select {
		case d.recvc <- buf[0:n]:
		
		case <-d.read_exit:
			return
		}
	}
}
