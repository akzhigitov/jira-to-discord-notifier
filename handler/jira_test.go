package handler

import "testing"

func TestParseDescriptionWithCode(t *testing.T)  {
	result := parseDescription("\r\n{code}\r\nselect\r\n{code}")

	if result != "\r\n```\r\nselect\r\n```\r\n"	{
		t.Error("not expected",result)
	}
}

func TestParseDescriptionWithCodeAndLanguage(t *testing.T)  {
	result := parseDescription("\r\n{code:java}\r\nselect\r\n{code}")

	if result != "\r\n```java\r\nselect\r\n```\r\n"	{
		t.Error("not expected",result)
	}
}

func TestParseDescriptionWithUnderline(t *testing.T)  {
	result := parseDescription(`+header+`)

	if result != "__header__"	{
		t.Error("not expected",result)
	}
}