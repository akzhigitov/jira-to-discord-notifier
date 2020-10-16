package handler

import "testing"

func TestParseDescription(t *testing.T)  {
	result := parseDescription(`
{code}
select
{code}`)

	if result != "\n```\nselect\n```\n"	{
		t.Error("not expected",result)
	}
}

func TestParseDescriptionAdvanced(t *testing.T)  {
	result := parseDescription(`
{code:java}
select
{code}`)

	if result != "\n```java\nselect\n```\n"	{
		t.Error("not expected",result)
	}
}