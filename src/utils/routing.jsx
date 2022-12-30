import React from 'react';
import { useLocation, useNavigate, useParams } from 'react-router-dom';

function withRouter(Component) {
  function ComponentWithRouterProp(props) {
    const location = useLocation();
    const navigate = useNavigate();
    const params = useParams();

    // eslint-disable-next-line react/jsx-props-no-spreading
    return <Component {...props} router={{ location, navigate, params }} />;
  }

  return ComponentWithRouterProp;
}

export default withRouter;
