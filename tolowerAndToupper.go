func tolower(c byte) byte {
  return c | 0x20
}

func toupper(c byte) byte {
	return c & ^byte(0x20)
}
