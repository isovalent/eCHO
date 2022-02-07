# Creating a new eCHO episode

## Resources

* Internal-only *eCHO-ing News* spreadsheet is [here](https://docs.google.com/spreadsheets/d/1Aq6sfOA06KbXNymt0UbBHYnvf-YWlR8e-oXjATzAigg/edit?usp=sharing)
  * Episode planning: host, guests, topic
  * Headlines
* eCHO Calendar (for anyone to subscribe to) is linked on the [README](README.md) 

## Preparation

**YouTube link and description**

* [ ] Set up the episode in StreamYard so that we have a YouTube URL, including the topic & guests
  * The same URL becomes the replay link for future viewers
* [ ] On YouTube, add episode to [eCHO playlist](https://www.youtube.com/playlist?list=PLDg_GiBbAx-mY3VFLPbLHcxo6wUjejAOC)
* [ ] On YouTube, set thumbnail image to the [default](https://github.com/isovalent/eCHO/blob/main/images/echo-cilium-ebpf-k8s.png) (or get one designed specific to the episode!) 
* [ ] Add URL & topic description to event in eCHO calendar 
* [ ] Add URL & topic description to event listing in the [README](README.md) 

**Episode notes**

* [ ] Draft the episode notes under a GitHub branch
  * There's a [tool](./episodes/README.md) that makes it easier to generate episode notes from the episode planning spreadsheet
* [ ] Sync to HackMD

**Invitations**

* [ ] Send invite to guest with Streamyard link (start 15 mins early)
* [ ] Invite a chat moderator to YouTube
* [ ] Share event on social media
* [ ] Post & pin link to #general channel in eBPF & Cilium Slack

**Content preparation**

Early in week of broadcast
* [ ] Agree episode outline with guest
* [ ] Prepare headlines (picking good ones from the eCHO-ing News sheet)
* [ ] Add links for the main topic: project repo, guest links, blogs, other resources...

## Episode ready

Just before broadcast
* [ ] Test sound & screen sharing with guest
* [ ] Reminders in social media and Slack

## After the episode

* [ ] Update & merge the episode notes as a PR
* [ ] Add link to the [README](README.md) and make sure next two episodes are listed
* [ ] Make sure episode is on the [YouTube playlist](https://www.youtube.com/playlist?list=PLDg_GiBbAx-mY3VFLPbLHcxo6wUjejAOC)
* [ ] Update thumbnail image

---------------------------------

# Archive instructions 

From when we used Streamlabs OBS:

**Graphics and streaming setup**

* [ ] Prepare any graphics e.g. a "lower third" for any guests
  * Use the [lower third template in Keynote](https://github.com/isovalent/eCHO/blob/main/streaming-assets/lower%20third%20example.key) 
  * Replace the lower third graphic, with the following animations:
    * Build in: Dissolve duration 1s
    * Build on: Dissolve duration 1s, start 5s after previous build
  * Export as a movie with transparency so it can be used as an overlay. See the export settings screenshot below
* [ ] Set up scenes in Streamlabs OBS
* [ ] Test sound & screen sharing with guest

**Keynote export settings for lower third overlay**

![](/images/overlay-keynote-export-settings.png)
