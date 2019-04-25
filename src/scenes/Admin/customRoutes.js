import React from 'react';
import { Route } from 'react-router-dom';
import AdminHome from './AdminHome';
import AdminUsers from './AdminUsers';

export default [
  <Route exact path="/system" component={AdminHome} />,
  <Route exact path="/system/users" component={AdminUsers} />,
];
