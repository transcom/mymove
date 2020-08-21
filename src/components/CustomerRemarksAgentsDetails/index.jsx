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
        <DataPoint columnHeaders={['Customer remarks']} dataRow={[customerRemarksBody]} />
      </div>
      <div className="container">
        <DataPoint columnHeaders={['Releasing agent']} dataRow={[releasingAgentBody]} />
      </div>
      <div className="container">
        <DataPoint columnHeaders={['Receiving agent']} dataRow={[receivingAgentBody]} />
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
