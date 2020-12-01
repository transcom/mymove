import React from 'react';
import { Route } from 'react-router-dom';
import UploadSearch from './Uploads/UploadSearch';
import UserSearch from './Users/UserSearch';

export default [
  <Route exact path="/uploads" component={UploadSearch} />,
  <Route exact path="/users" component={UserSearch} />,
];
