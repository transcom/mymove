import React, { Component } from 'react';
import PropTypes from 'prop-types';

import { HHGDetailsForm } from 'components/Customer/HHGDetailsForm';

// eslint-disable-next-line react/prefer-stateless-function
export class HHGMoveSetup extends Component {
  render() {
    const { pageList, pageKey } = this.props;
    return (
      <div>
        <h3>Now lets arrange details for the professional movers</h3>
        <HHGDetailsForm pageList={pageList} pageKey={pageKey} />
      </div>
    );
  }
}

HHGMoveSetup.propTypes = {
  pageList: PropTypes.arrayOf(PropTypes.string),
  pageKey: PropTypes.string,
};

HHGMoveSetup.defaultProps = {
  pageList: [],
  pageKey: '',
};

export default HHGMoveSetup;
