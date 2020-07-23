import React from 'react';
import { arrayOf, string, shape, bool } from 'prop-types';

import HHGDetailsForm from 'components/Customer/HHGDetailsForm';

const HHGMoveSetup = ({ pageList, pageKey, match }) => (
  <div>
    <h3>Now lets arrange details for the professional movers</h3>
    <HHGDetailsForm pageList={pageList} pageKey={pageKey} match={match} />
  </div>
);

HHGMoveSetup.propTypes = {
  pageList: arrayOf(string).isRequired,
  pageKey: string.isRequired,
  match: shape({
    isExact: bool.isRequired,
    params: shape({
      moveId: string.isRequired,
    }),
    path: string.isRequired,
    url: string.isRequired,
  }).isRequired,
};

export default HHGMoveSetup;
