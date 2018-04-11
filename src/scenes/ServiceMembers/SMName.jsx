import React, { Component } from 'react';
import NameForm from './NameForm';
import { connect } from 'react-redux';
import { bindActionCreators } from 'redux';
import PropTypes from 'prop-types';

import { setPendingSMNameData, loadServiceMember } from './ducks';

class SMName extends Component {
  componentDidMount() {
    document.title = 'Transcom PPP: Service Member Name';
  }
  onNameDataEntry = values => {
    this.props.setPendingSMNameData(values);
  };
  render() {
    // const { pendingSMNameData, currentServiceMember } = this.props;
    return (
      <div>
        <NameForm onSubmit={this.onNameDataEntry} />
      </div>
    );
  }
}

SMName.propTypes = {
  currentServiceMember: PropTypes.object,
  pendingSMNameData: PropTypes.object,
  setPendingSMNameData: PropTypes.func.isRequired,
};

function mapDispatchToProps(dispatch) {
  return bindActionCreators(
    { setPendingSMNameData, loadServiceMember },
    dispatch,
  );
}

function mapStateToProps(state) {
  return { ...state.serviceMember };
}

export default connect(mapDispatchToProps, mapDispatchToProps)(SMName);
