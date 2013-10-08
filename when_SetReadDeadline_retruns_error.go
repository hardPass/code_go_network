// SetReadDeadline implements the Conn SetReadDeadline method.
func (c *conn) SetReadDeadline(t time.Time) error {
	if !c.ok() {
		return syscall.EINVAL
	}
	return setReadDeadline(c.fd, t)
}

func (c *conn) ok() bool { return c != nil && c.fd != nil }


func setWriteDeadline(fd *netFD, t time.Time) error {
	return setDeadlineImpl(fd, t, 'w')
}

func setDeadline(fd *netFD, t time.Time) error {
	return setDeadlineImpl(fd, t, 'r'+'w')
}

func setDeadlineImpl(fd *netFD, t time.Time, mode int) error {
	d := t.UnixNano()
	if t.IsZero() {
		d = 0
	}
	if err := fd.incref(false); err != nil {
		return err
	}
	runtime_pollSetDeadline(fd.pd.runtimeCtx, d, mode)
	fd.decref()
	return nil
}

// Add a reference to this fd.
// If closing==true, pollDesc must be locked; mark the fd as closing.
// Returns an error if the fd cannot be used.
func (fd *netFD) incref(closing bool) error {
	fd.sysmu.Lock()
	if fd.closing {
		fd.sysmu.Unlock()
		return errClosing
	}
	fd.sysref++
	if closing {
		fd.closing = true
	}
	fd.sysmu.Unlock()
	return nil
}
