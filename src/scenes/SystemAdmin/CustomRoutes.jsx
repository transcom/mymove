import React from 'react';
import { Route } from 'react-router-dom-old';
import UploadSearch from './Uploads/UploadSearch';

const routes = [
  //Custom route for search by id for uploads
  <Route exact path="/uploads" component={UploadSearch} />,
];

export default routes;
