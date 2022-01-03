import { get } from 'lodash';
import PropTypes from 'prop-types';
import React, { Component } from 'react';
import { connect } from 'react-redux';
import { bindActionCreators } from 'redux';
import {
  loadDutyLocationTransportationOffice,
  selectDutyLocationTransportationOffice,
} from 'shared/Entities/modules/transportationOffices';

export class TransportationOfficeContactInfo extends Component {
  componentDidMount() {
    const { dutyLocation } = this.props;
    this.props.loadDutyLocationTransportationOffice(dutyLocation.id);
  }
  render() {
    const { isOrigin, dutyLocation, transportationOffice } = this.props;
    const transportationOfficeName = get(transportationOffice, 'name');
    const officeName =
      transportationOfficeName && get(dutyLocation, 'name') !== transportationOfficeName
        ? transportationOffice.name
        : 'Transportation Office';
    const contactInfo = Boolean(get(transportationOffice, 'phone_lines[0]'));
    return (
      <div className="titled_block">
        {dutyLocation && <strong>{dutyLocation.name}</strong>}
        <div>
          {isOrigin ? 'Origin' : 'Destination'} {officeName}
        </div>
        <div>{contactInfo ? get(transportationOffice, 'phone_lines[0]') : 'Contact Info Not Available'}</div>
      </div>
    );
  }
}
TransportationOfficeContactInfo.propTypes = {
  loadDutyLocationTransportationOffice: PropTypes.func.isRequired,
  dutyLocation: PropTypes.shape({
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
    transportationOffice: selectDutyLocationTransportationOffice(state),
  };
};

function mapDispatchToProps(dispatch) {
  return bindActionCreators({ loadDutyLocationTransportationOffice }, dispatch);
}

export default connect(mapStateToProps, mapDispatchToProps)(TransportationOfficeContactInfo);
