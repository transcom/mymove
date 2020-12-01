import { get } from 'lodash';
import PropTypes from 'prop-types';
import React, { Component } from 'react';
import { connect } from 'react-redux';
import { bindActionCreators } from 'redux';
import {
  loadDutyStationTransportationOffice,
  selectDutyStationTransportationOffice,
} from 'shared/Entities/modules/transportationOffices';

export class TransportationOfficeContactInfo extends Component {
  componentDidMount() {
    const { dutyStation } = this.props;
    this.props.loadDutyStationTransportationOffice(dutyStation.id);
  }
  render() {
    const { isOrigin, dutyStation, transportationOffice } = this.props;
    const transportationOfficeName = get(transportationOffice, 'name');
    const officeName =
      transportationOfficeName && get(dutyStation, 'name') !== transportationOfficeName
        ? transportationOffice.name
        : 'Transportation Office';
    const contactInfo = Boolean(get(transportationOffice, 'phone_lines[0]'));
    return (
      <div className="titled_block">
        {dutyStation && <strong>{dutyStation.name}</strong>}
        <div>
          {isOrigin ? 'Origin' : 'Destination'} {officeName}
        </div>
        <div>{contactInfo ? get(transportationOffice, 'phone_lines[0]') : 'Contact Info Not Available'}</div>
      </div>
    );
  }
}
TransportationOfficeContactInfo.propTypes = {
  loadDutyStationTransportationOffice: PropTypes.func.isRequired,
  dutyStation: PropTypes.shape({
    name: PropTypes.string.isRequired,
    id: PropTypes.string.isRequired,
    transportation_office: PropTypes.shape({ phone_lines: PropTypes.array }),
  }),
  isOrigin: PropTypes.bool.isRequired,
};
TransportationOfficeContactInfo.defaultProps = {
  transportationOffice: {},
  isOrigin: false,
};

const mapStateToProps = (state, ownProps) => {
  return {
    transportationOffice: selectDutyStationTransportationOffice(state),
  };
};

function mapDispatchToProps(dispatch) {
  return bindActionCreators({ loadDutyStationTransportationOffice }, dispatch);
}

export default connect(mapStateToProps, mapDispatchToProps)(TransportationOfficeContactInfo);
