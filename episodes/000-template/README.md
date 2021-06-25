# Episode {{.Number}}: {{.Title}}

[YouTube]({{.YouTube}})
{{friendly .Date}}

with [{{.GuestName}}]({{.GuestURL}}), hosted by [{{.HostName}}]({{.HostURL}})

## Headlines

## Links

## {{.Title}}


# Updating the main README episode list

1. Remove this episode from "Upcoming episodes"
2. Add a new episode in two weeks' time to "Upcoming episodes"

- Episode #{{upcoming .Number}}: TBD
{{if .Eu}}  - 6am PT, 9am ET, 2pm UK, 3pm Central Europe{{else}}  - 11am PT, 2pm ET, 7pm UK, 8pm Central Europe{{end}} - {{friendly (upcomingDate .Date)}}
  - [Convert to your timezone / get calendar link](https://www.timeanddate.com/worldclock/fixedtime.html?msg=eBPF+%26+Cilium+Office+Hours&iso={{iso (upcomingDate .Date)}}T{{if .Eu}}14{{else}}19{{end}}&p1=136&am=30)
  - [Link to YouTube TBD]()

3. Add this episode to the history

- Episode #{{.Number}}: [{{.Title}}]({{.YouTube}}) with [{{.GuestName}}]({{.GuestURL}})
  - [Show notes]({{.ShowNotesURL}})


