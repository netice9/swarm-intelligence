import React, { Component } from 'react'
import { Terminal } from 'xterm'
import "xterm/dist/xterm.css"

export default class LogTailer extends Component {



  componentDidMount() {
      this.term = new Terminal()
      this.term.open(this.container)
      if (this.props.url) {
        const decoder = new TextDecoder('utf-8')
        fetch(this.props.url, {credentials: 'same-origin'})
        .then((response) => response.body)
        .then((body) => body.getReader())
        .then((reader) => {
          this.reader = reader
          const readChunk = ({ done, value }) => {
            this.term.write(decoder.decode(value).replace(/\r?\n/g, "\r\n"))
            if (done) {
              this.term.write("closed ...")
              return
            }
            return reader.read().then(readChunk);
          }
          reader.read().then(readChunk)
        })
      }


  }

  componentWillUnmount() {
    if (this.reader) {
      this.reader.cancel()
    }
    if (this.xterm) {
      this.xterm.destroy()
      this.xterm = null
    }
  }



  render() {
    return <div ref={ref => (this.container = ref)} />;
  }
}
