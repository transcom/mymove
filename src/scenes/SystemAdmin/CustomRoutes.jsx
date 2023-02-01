import React from 'react';
import { Route } from 'react-router-dom';
import UploadSearch from './Uploads/UploadSearch';
import { CustomRoutes } from 'react-admin';

const Routes = () => (
  <CustomRoutes>
    {/* Custom route for search by id for uploads */}
    <Route exact path="/uploads" component={UploadSearch} />
  </CustomRoutes>
);

export default Routes;
