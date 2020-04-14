import React from 'react';
import { GovBanner } from '@trussworks/react-uswds';
import { withRouter } from 'react-router-dom';
import { connect } from 'react-redux';

const Accessorials = () => {
  return (
    <>
      <GovBanner />
      <h1>This is where we will put our accessorial components!</h1>
    </>
  );
};

export default withRouter(connect()(Accessorials));
