import React from 'react';
import { Route } from 'react-router-dom';
import AdminHome from './AdminHome';
import AdminUsers from './AdminUsers';

export default [
  <Route exact path="/portal" component={AdminHome} />,
  <Route exact path="/portal/users" component={AdminUsers} />,
];
