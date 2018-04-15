import React from 'react';
import ReactDOM from 'react-dom';
import {BrowserRouter, Switch, Route} from 'react-router-dom'
import '../style/main.css'
import '../style/main.scss'
import { Provider } from 'react-redux'
import { createStore, applyMiddleware } from 'redux'
import reducers from './reducers'
import { swarmStateUpdate, websocketConnected, websocketDisconnected } from './actions'
import WebsocketSpinner from './components/websocket_spinner'

import Namespaces from './components/namespaces'
import DeployStack from './components/deploy_stack'
import NavBar from './components/nav_bar'
import ManageCredentials from './components/manage_credentials'

const store = createStore(reducers)

ReactDOM.render(
  <Provider store={store}>
    <WebsocketSpinner>
      <BrowserRouter>
        <div>
          <NavBar />
          <Route exact path="/" component={Namespaces} />
          <Route exact path="/deploy_stack" component={DeployStack} />
          <Route exact path="/manage_credentials" component={ManageCredentials} />
        </div>
      </BrowserRouter>
    </WebsocketSpinner>
  </Provider>,
  document.getElementById('root')
)

const protocol = location.protocol === "https:" ? "wss" : "ws"
const url = `${protocol}://${location.host}/api/state`

const connectToWebsocket = () => {
  var socket = new WebSocket(url)

  socket.onopen = (evt) => {
    store.dispatch(websocketConnected())
  }

  socket.onclose = (evt) => {
    store.dispatch(websocketDisconnected())
    setTimeout(connectToWebsocket, 1000)
  }

  socket.onmessage = (event) => {
    store.dispatch(swarmStateUpdate(JSON.parse(event.data)))
  }
}

connectToWebsocket()
