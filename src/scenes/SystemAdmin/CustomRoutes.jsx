import React from 'react';
import { Route } from 'react-router-dom';
import UploadSearch from './Uploads/UploadSearch';
import GexSearch from './Gex/GexSearch';

export default [
  <Route exact path="/uploads" component={UploadSearch} />,
  <Route exact path="/gex" component={GexSearch} />,
];
