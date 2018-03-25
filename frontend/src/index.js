import React from 'react';
import ReactDOM from 'react-dom';
import {BrowserRouter, Switch, Route} from 'react-router-dom'
import '../style/main.css'
import '../style/main.scss'
import Index from './components/index'
import NavBar from './components/nav_bar'

import { Provider } from 'react-redux'
import { createStore, applyMiddleware } from 'redux'
import reducers from './reducers'
import promise from 'redux-promise'


const createStoreWithMiddleware = applyMiddleware(promise)(createStore)

ReactDOM.render(
  <Provider store={createStoreWithMiddleware(reducers)}>
    <div>
      <BrowserRouter>
        <div>
          <NavBar />
          <Route exact path="/" component={Index} />
        </div>
      </BrowserRouter>
    </div>
  </Provider>,
  document.getElementById('root')
);
