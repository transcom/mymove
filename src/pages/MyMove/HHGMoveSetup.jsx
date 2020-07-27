import React from 'react';
import { arrayOf, string, shape, bool, func } from 'prop-types';

import HHGDetailsForm from 'components/Customer/HHGDetailsForm';

const HHGMoveSetup = ({ pageList, pageKey, match, push }) => (
  <div>
    <h3>Now lets arrange details for the professional movers</h3>
    <HHGDetailsForm pageList={pageList} pageKey={pageKey} match={match} push={push} />
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
  push: func.isRequired,
};

export default HHGMoveSetup;
