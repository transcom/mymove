import React from 'react';
import propTypes from 'prop-types';
import { get } from 'lodash';

import DataPoint from 'components/DataPoint';

const CustomerRemarksAgentsDetails = ({ customerRemarks, releasingAgent, receivingAgent }) => {
  const customerRemarksBody = <>{customerRemarks}</>;
  const releasingAgentBody = (
    <>
      {(get(releasingAgent, 'firstName') || get(releasingAgent, 'lastName')) && (
        <>
          {`${get(releasingAgent, 'firstName')} ${get(releasingAgent, 'lastName')}`}
          <br />
        </>
      )}

      {get(releasingAgent, 'phone') && (
        <>
          {get(releasingAgent, 'phone')}
          <br />
        </>
      )}
      {get(releasingAgent, 'email')}
    </>
  );
  const receivingAgentBody = (
    <>
      {(get(receivingAgent, 'firstName') || get(receivingAgent, 'lastName')) && (
        <>
          {`${get(receivingAgent, 'firstName')} ${get(receivingAgent, 'lastName')}`}
          <br />
        </>
      )}

      {get(receivingAgent, 'phone') && (
        <>
          {get(receivingAgent, 'phone')}
          <br />
        </>
      )}
      {get(receivingAgent, 'email')}
    </>
  );

  return (
    <>
      <div className="container">
        <DataPoint header="Customer remarks" body={customerRemarksBody} />
      </div>
      <div className="container">
        <DataPoint header="Releasing agent" body={releasingAgentBody} />
      </div>
      <div className="container">
        <DataPoint header="Receiving agent" body={receivingAgentBody} />
      </div>
    </>
  );
};

CustomerRemarksAgentsDetails.propTypes = {
  customerRemarks: propTypes.string.isRequired,
  releasingAgent: propTypes.shape({
    firstName: propTypes.string,
    lastName: propTypes.string,
    phone: propTypes.string,
    email: propTypes.string,
  }).isRequired,
  receivingAgent: propTypes.shape({
    firstName: propTypes.string,
    lastName: propTypes.string,
    phone: propTypes.string,
    email: propTypes.string,
  }).isRequired,
};

export default CustomerRemarksAgentsDetails;
