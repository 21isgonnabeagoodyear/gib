benis :DDDD
# Design Doc
## General Concepts
- Try and keep dependencies as low as possible
  - Obviously don't reinvent the wheel if a package does something really well, but as a rule of thumb, try and stay standard-library
  - Try and stick to plain `net/http`; we'll mess with gorilla if it appears to become necessary
  - As a corrolary, packages that we probably *will* need are:
    - `github.com/lib/pq` (just do postgres for now; we can add in sqlite and stuff later
    - drivers for `ffmpeg` and/or `imagemagick` for thumbs
- Just stick to http; don't try and provide cgi, fast or traditional
- Don't bother serving https either
  - `crypto/tls` is vulnerable to timing atacks
  - `github.com/spacemonkeygo/openssl` introduces a cgo dependency that doesn't appear to work on any of yuuko's systems
  - Leave it to reverse-proxied forward-facing httpds (nginx, apache, etc.) for separation of tasks
- thumbnailers should be easily "pluggable" by filetype; some ideas:
  - Normal thumbnails for images (preserving PNG transparency)
  - Typical thumbnails for video
  - spectrographs for audio
  - first page for pdf/cbr
- board settings (max upload size, etc.) should all be customizable, with defaults inherited from global config
- config is in either json or yaml
- traditional boards model
