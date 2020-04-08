import React from 'react';
import { GovBanner } from '@trussworks/react-uswds';
import { withRouter } from 'react-router-dom';
import { connect } from 'react-redux';

const TOO = () => {
  return (
    <>
      <GovBanner />
      <h1>This is where we will put our accessorial components!</h1>
    </>
  );
};

const mapStateToProps = () => {
  return {};
};
const mapDispatchToProps = {};

export default withRouter(connect(mapStateToProps, mapDispatchToProps)(TOO));
