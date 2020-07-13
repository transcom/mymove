import React, { Component } from 'react';
import PropTypes from 'prop-types';

import { HHGDetailsForm } from 'components/Customer/HHGDetailsForm';

// eslint-disable-next-line react/prefer-stateless-function
export class HHGMoveSetup extends Component {
  render() {
    const { pages, pageKey } = this.props;
    return (
      <div>
        <h3>Now lets arrange details for the professional movers</h3>
        <HHGDetailsForm pages={pages} pageKey={pageKey} />
      </div>
    );
  }
}

HHGMoveSetup.propTypes = {
  pages: PropTypes.arrayOf(PropTypes.shape({})).isRequired,
  pageKey: PropTypes.string.isRequired,
};

export default HHGMoveSetup;
