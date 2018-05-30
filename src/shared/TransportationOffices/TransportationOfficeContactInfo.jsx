import { get } from 'lodash';
import PropTypes from 'prop-types';
import React, { Component } from 'react';
import { connect } from 'react-redux';
import { bindActionCreators } from 'redux';
import {
  loadDutyStationTransportationOffice,
  getDutyStationTransportationOffice,
} from 'shared/TransportationOffices/ducks';
import { no_op } from 'shared/utils';

export class TransportationOfficeContactInfo extends Component {
  componentDidMount() {
    this.props.loadDutyStationTransportationOffice(
      get(this.props, 'dutyStation.id'),
    );
  }
  render() {
    const { isOrigin, dutyStation, transportationOffice } = this.props;
    const transportationOfficeName = get(transportationOffice, 'name');
    const officeName =
      transportationOfficeName &&
      get(dutyStation, 'name') !== transportationOfficeName
        ? transportationOffice.name
        : 'Transportation Office';
    const contactInfo = Boolean(get(transportationOffice, 'phone_lines[0]'));
    return (
      <div className="titled_block">
        {dutyStation && <strong>{dutyStation.name}</strong>}
        <div>
          {isOrigin ? 'Origin' : 'Destination'} {officeName}
        </div>
        <div>
          {contactInfo
            ? get(transportationOffice, 'phone_lines[0]')
            : 'Contact Info Not Available'}
        </div>
      </div>
    );
  }
}
TransportationOfficeContactInfo.propTypes = {
  getDutyStationTransportationOffice: PropTypes.func.isRequired,
  loadDutyStationTransportationOffice: PropTypes.func.isRequired,
  dutyStation: PropTypes.shape({
    name: PropTypes.string.isRequired,
    id: PropTypes.string.isRequired,
  }),
  isOrigin: PropTypes.bool.isRequired,
};
TransportationOfficeContactInfo.defaultProps = {
  getDutyStationTransportationOffice: no_op,
  loadDutyStationTransportationOffice: no_op,
  isOrigin: false,
};

const mapStateToProps = (state, ownProps) => ({
  transportationOffice: getDutyStationTransportationOffice(
    state,
    get(ownProps, 'dutyStation.id'),
  ),
});

function mapDispatchToProps(dispatch) {
  return bindActionCreators({ loadDutyStationTransportationOffice }, dispatch);
}

export default connect(mapStateToProps, mapDispatchToProps)(
  TransportationOfficeContactInfo,
);
