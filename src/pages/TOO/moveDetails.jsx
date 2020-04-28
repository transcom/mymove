import React, { Component } from 'react';
import { withRouter } from 'react-router-dom';
import { connect } from 'react-redux';
import '../../index.scss';
import '../../ghc_index.scss';

class MoveDetails extends Component {
  componentDidMount() {}

  render() {
    return (
      <div className="maxw-desktop-lg" data-cy="too-move-details">
        <h1>Move Details</h1>
      </div>
    );
  }
}

const mapStateToProps = () => {
  return {};
};

const mapDispatchToProps = {};

export default withRouter(connect(mapStateToProps, mapDispatchToProps)(MoveDetails));
