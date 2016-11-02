/*@flow*/
import React from 'react';
import { Router, Route } from 'dva/router';
import dva from 'dva';
import App from './App';
import model from './model';
import './index.css';

const app = dva();
app.model(model);

app.router(({ history }) =>
    <Router history={history}>
    <Route path="/" component={App} />
  </Router>
);
app.start('#root');
