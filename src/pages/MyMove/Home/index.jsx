/* eslint-disable react/no-unused-prop-types */
/* eslint-disable no-console */
import React, { Component } from 'react';
import { func, arrayOf, shape, string, node } from 'prop-types';
import { connect } from 'react-redux';

import styles from './Home.module.scss';

import Helper from 'components/Customer/Home/Helper';
import Step from 'components/Customer/Home/Step';
import DocsUploaded from 'components/Customer/Home/DocsUploaded';
import ShipmentList from 'components/Customer/Home/ShipmentList';
import Contact from 'components/Customer/Home/Contact';
import { showLoggedInUser as showLoggedInUserAction } from 'shared/Entities/modules/user';
import { selectServiceMemberFromLoggedInUser } from 'shared/Entities/modules/serviceMembers';
import { selectUploadedOrders, selectActiveOrLatestOrdersFromEntities } from 'shared/Entities/modules/orders';
import { selectActiveOrLatestMove } from 'shared/Entities/modules/moves';
import {
  selectMTOShipmentsByMoveId,
  loadMTOShipments as loadMTOShipmentsAction,
} from 'shared/Entities/modules/mtoShipments';

const Description = ({ children }) => <p className={styles.description}>{children}</p>;

Description.propTypes = {
  children: node.isRequired,
};

class Home extends Component {
  componentDidMount() {
    const { showLoggedInUser, loadMTOShipments, move } = this.props;
    showLoggedInUser();
    loadMTOShipments(move.id);
  }

  get hasOrders() {
    const { orders, uploadedOrderDocuments } = this.props;
    return !!Object.keys(orders).length && !!uploadedOrderDocuments.length;
  }

  get hasShipments() {
    const { shipments } = this.props;
    return this.hasOrders && !!shipments.length;
  }

  get getHelperHeaderText() {
    if (!this.hasOrders) {
      return 'Next step: Add your orders';
    }

    if (!this.hasShipments) {
      return 'Gather this info, then plan your shipments';
    }

    if (this.hasShipments) {
      return 'Time to submit your move';
    }

    return '';
  }

  renderHelperListItems = (helperList) => {
    return helperList.map((listItemText) => (
      <li key={listItemText}>
        <span>{listItemText}</span>
      </li>
    ));
  };

  renderHelperDescription = () => {
    if (!this.hasOrders) {
      return (
        <ul>
          {this.renderHelperListItems([
            'If you have a hard copy, you can take photos of each page',
            'If you have a PDF, you can upload that',
          ])}
        </ul>
      );
    }

    if (!this.hasShipments) {
      return (
        <ul>
          {this.renderHelperListItems([
            'Preferred moving details',
            'Destination address (your new place, your duty station ZIP, or somewhere else',
            'Names and contact info for anyone you authorize to act on your behalf',
          ])}
        </ul>
      );
    }

    if (this.hasShipments) {
      return (
        <ul>
          {this.renderHelperListItems([
            "Double check the info you've entered",
            'Sign the legal agreement',
            "You'll hear from a move counselor or your transportation office within a few days",
          ])}
        </ul>
      );
    }

    return null;
  };

  handleShipmentClick = (shipment) => {
    console.log('this is the shipment', shipment);
  };

  render() {
    const { serviceMember, uploadedOrderDocuments, shipments } = this.props;
    return (
      <div className={`usa-prose grid-container ${styles['grid-container']}`}>
        <header className={styles['customer-header']}>
          <h2>
            {serviceMember.first_name} {serviceMember.last_name}
          </h2>
          <p>
            You&apos;re leaving <strong>{serviceMember.current_station.name}</strong>
          </p>
        </header>
        <Helper title={this.getHelperHeaderText}>{this.renderHelperDescription()}</Helper>
        <Step
          complete={serviceMember.is_profile_complete}
          completedHeaderText="Profile complete"
          editBtnLabel="Edit"
          headerText="Profile complete"
          step="1"
          onEditClick={(e) => {
            e.preventDefault();
            console.log('edit clicked');
          }}
        >
          <Description>Make sure to keep your personal information up to date during your move</Description>
        </Step>

        <Step
          actionBtnLabel={!this.hasOrders ? 'Add orders' : ''}
          complete={this.hasOrders}
          completedHeaderText="Orders uploaded"
          editBtnLabel={this.hasOrders ? 'Edit' : ''}
          onEditClick={() => console.log('edit button clicked')}
          headerText="Upload orders"
          onActionBtnClick={() => console.log('some action')}
          step="2"
        >
          {this.hasOrders ? (
            <DocsUploaded files={uploadedOrderDocuments} />
          ) : (
            <Description>Upload photos of each page, or upload a PDF.</Description>
          )}
        </Step>

        <Step
          actionBtnDisabled={!this.hasOrders}
          actionBtnLabel={this.hasShipments ? 'Add another shipment' : 'Plan your shipments'}
          complete={this.hasShipments}
          completedHeaderText="Shipments"
          headerText="Shipment selection"
          secondaryBtn={this.hasShipments}
          secondaryClassName="margin-top-2"
          step="3"
        >
          {this.hasShipments ? (
            <ShipmentList shipments={shipments} onShipmentClick={this.handleShipmentClick} />
          ) : (
            <Description>
              Tell us where you&apos;re going and when you want to get there. We&apos;ll help you set up shipments to
              make it work.
            </Description>
          )}
        </Step>

        <Step
          actionBtnDisabled={!this.hasShipments}
          actionBtnLabel="Review and submit"
          containerClassName="margin-bottom-8"
          headerText="Confirm move request"
          onActionBtnClick={() => console.log('some action')}
          step="4"
        >
          <Description>
            Review your move details and sign the legal paperwork, then send the info on to your move counselor
          </Description>
        </Step>
        <Contact
          header="Contacts"
          dutyStationName="Seymour Johnson AFB"
          officeType="Origin Transportation Office"
          telephone="(919) 722-5458"
        />
      </div>
    );
  }
}

Home.propTypes = {
  orders: shape({}).isRequired,
  serviceMember: shape({
    first_name: string.isRequired,
    last_name: string.isRequired,
  }).isRequired,
  showLoggedInUser: func.isRequired,
  loadMTOShipments: func.isRequired,
  shipments: arrayOf(
    shape({
      id: string,
      locator: string,
      select_move_type: string,
    }),
  ).isRequired,
  uploadedOrderDocuments: arrayOf(
    shape({
      filename: string.isRequired,
    }),
  ).isRequired,
  move: shape({}).isRequired,
};

const mapStateToProps = (state) => {
  const serviceMember = selectServiceMemberFromLoggedInUser(state);
  const move = selectActiveOrLatestMove(state);
  console.log('move', move);
  return {
    orders: selectActiveOrLatestOrdersFromEntities(state),
    uploadedOrderDocuments: selectUploadedOrders(state),
    serviceMember,
    // TODO: change when we support PPM shipments as well
    shipments: selectMTOShipmentsByMoveId(state, move.id),
    move,
  };
};

// in order to avoid setting up proxy server only for storybook, pass in stub function so API requests don't fail
const mergeProps = (stateProps, dispatchProps, ownProps) => ({
  ...stateProps,
  ...dispatchProps,
  ...ownProps,
});

const mapDispatchToProps = {
  showLoggedInUser: showLoggedInUserAction,
  loadMTOShipments: loadMTOShipmentsAction,
};

export default connect(mapStateToProps, mapDispatchToProps, mergeProps)(Home);
