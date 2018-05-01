import PropTypes from 'prop-types';
import React, { Component } from 'react';
import { connect } from 'react-redux';
import { bindActionCreators } from 'redux';

import { updateServiceMember, loadServiceMember } from './ducks';
import { no_op } from 'shared/utils';
import WizardPage from 'shared/WizardPage';

import DutyStationSearchBox from 'scenes/ServiceMembers/DutyStationSearchBox';

import './DutyStation.css';

export class DutyStation extends Component {
  constructor(props) {
    super(props);

    this.state = {
      value: null,
    };
    this.stationOnChange = this.stationOnChange.bind(this);
  }

  stationOnChange = newStation => {
    this.setState({ value: newStation });
  };
  componentDidMount() {
    this.props.loadServiceMember(this.props.match.params.serviceMemberId);
  }

  handleSubmit = () => {
    if (this.state.value) {
      this.props.updateServiceMember({ current_station: this.state.value });
    }
  };

  render() {
    const { pages, pageKey, hasSubmitSuccess, error } = this.props;
    // TODO: make sure isvalid is accurate
    const isValid = !!this.state.value;
    return (
      <WizardPage
        handleSubmit={this.handleSubmit}
        isAsync={true}
        pageList={pages}
        pageKey={pageKey}
        pageIsValid={isValid}
        hasSucceeded={hasSubmitSuccess}
        error={error}
      >
        <form className="duty-station" onSubmit={no_op}>
          <h1 className="sm-heading">Current Duty Station</h1>
          <DutyStationSearchBox onChange={this.stationOnChange} />
        </form>
      </WizardPage>
    );
  }
}
DutyStation.propTypes = {
  error: PropTypes.object,
  hasSubmitSuccess: PropTypes.bool.isRequired,
  updateServiceMember: PropTypes.func.isRequired,
};

function mapDispatchToProps(dispatch) {
  return bindActionCreators(
    { updateServiceMember, loadServiceMember },
    dispatch,
  );
}
function mapStateToProps(state) {
  const props = {
    ...state.serviceMember,
  };
  return props;
}
export default connect(mapStateToProps, mapDispatchToProps)(DutyStation);
