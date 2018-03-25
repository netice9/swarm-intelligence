import React from 'react';
import ReactDOM from 'react-dom';
import {BrowserRouter, Switch, Route} from 'react-router-dom'
import '../style/main.css'
import '../style/main.scss'
import Index from './components/index'
import DeployStack from './components/deploy_stack'
import NavBar from './components/nav_bar'
import { Provider } from 'react-redux'
import { createStore, applyMiddleware } from 'redux'
import reducers from './reducers'
import { swarmStateUpdate } from './actions'
const store = createStore(reducers)

ReactDOM.render(
  <Provider store={store}>
    <div>
      <BrowserRouter>
        <div>
          <NavBar />
          <Route exact path="/" component={Index} />
          <Route exact path="/deploy_stack" component={DeployStack} />
        </div>
      </BrowserRouter>
    </div>
  </Provider>,
  document.getElementById('root')
)

const protocol = location.protocol === "https:" ? "wss" : "ws"
var socket = new WebSocket(`${protocol}://${location.host}/api/state`)
console.log(socket)
socket.onmessage = (event) => {
  store.dispatch(swarmStateUpdate(JSON.parse(event.data)))
}
