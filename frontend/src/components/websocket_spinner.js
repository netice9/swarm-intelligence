import React, { Component } from 'react'
import { connect } from 'react-redux'
import Loadable from 'react-loading-overlay'

class WebsocketSpinner extends Component {
  render() {
    return (
      <Loadable
        spinner
        active={!this.props.websocketConnected}
        text="Connecting to server"
      >
        { this.props.children }
      </Loadable>
    )
  }
}

const mapStateToProps = (state) => {
  return {
    websocketConnected: state.websocketConnected
  }
}

export default connect(mapStateToProps)(WebsocketSpinner)
