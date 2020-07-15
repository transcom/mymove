import React, { Component } from 'react';
import PropTypes from 'prop-types';

import HHGDetailsFormComponent from 'components/Customer/HHGDetailsForm';

// eslint-disable-next-line react/prefer-stateless-function
export class HHGMoveSetup extends Component {
  render() {
    const { pageList, pageKey } = this.props;
    return (
      <div>
        <h3>Now lets arrange details for the professional movers</h3>
        <HHGDetailsFormComponent pageList={pageList} pageKey={pageKey} />
      </div>
    );
  }
}

HHGMoveSetup.propTypes = {
  pageList: PropTypes.arrayOf(PropTypes.string).isRequired,
  pageKey: PropTypes.string.isRequired,
};

export default HHGMoveSetup;
