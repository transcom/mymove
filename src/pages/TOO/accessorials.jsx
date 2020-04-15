import React from 'react';
import { withRouter } from 'react-router-dom';
import { connect } from 'react-redux';

const Accessorials = () => {
  return (
    <>
      <h1>This is where we will put our accessorial components!</h1>
    </>
  );
};

export default withRouter(connect()(Accessorials));
