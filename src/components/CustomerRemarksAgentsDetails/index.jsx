import React from 'react';
import propTypes from 'prop-types';
import { get } from 'lodash';

import DataTable from 'components/DataTable';

const CustomerRemarksAgentsDetails = ({ customerRemarks, releasingAgent, receivingAgent }) => {
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
        <DataTable columnHeaders={['Customer remarks']} dataRow={[customerRemarks]} />
      </div>
      <div className="container">
        <DataTable columnHeaders={['Releasing agent']} dataRow={[releasingAgentBody]} />
      </div>
      <div className="container">
        <DataTable columnHeaders={['Receiving agent']} dataRow={[receivingAgentBody]} />
      </div>
    </>
  );
};

CustomerRemarksAgentsDetails.propTypes = {
  customerRemarks: propTypes.string,
  releasingAgent: propTypes.shape({
    firstName: propTypes.string,
    lastName: propTypes.string,
    phone: propTypes.string,
    email: propTypes.string,
  }),
  receivingAgent: propTypes.shape({
    firstName: propTypes.string,
    lastName: propTypes.string,
    phone: propTypes.string,
    email: propTypes.string,
  }),
};

CustomerRemarksAgentsDetails.defaultProps = {
  customerRemarks: '',
  releasingAgent: {},
  receivingAgent: {},
};

export default CustomerRemarksAgentsDetails;
