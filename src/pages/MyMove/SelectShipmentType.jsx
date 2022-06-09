/* eslint-disable react/jsx-props-no-spreading */
import React, { Component } from 'react';
import { connect } from 'react-redux';
import { func, arrayOf } from 'prop-types';
import { GridContainer, Grid, Alert } from '@trussworks/react-uswds';
import { generatePath } from 'react-router';

import ConnectedMoveInfoModal from 'components/Customer/modals/MoveInfoModal/MoveInfoModal';
import ConnectedStorageInfoModal from 'components/Customer/modals/StorageInfoModal/StorageInfoModal';
import SelectableCard from 'components/Customer/SelectableCard';
import WizardNavigation from 'components/Customer/WizardNavigation/WizardNavigation';
import ScrollToTop from 'components/ScrollToTop';
import { generalRoutes, customerRoutes } from 'constants/routes';
import styles from 'pages/MyMove/SelectShipmentType.module.scss';
import { patchMove, getResponseError } from 'services/internalApi';
import { SHIPMENT_OPTIONS } from 'shared/constants';
import { loadMTOShipments as loadMTOShipmentsAction } from 'shared/Entities/modules/mtoShipments';
import { updateMove as updateMoveAction } from 'store/entities/actions';
import { selectCurrentMove, selectMTOShipmentsForCurrentMove } from 'store/entities/selectors';
import formStyles from 'styles/form.module.scss';
import { MoveTaskOrderShape } from 'types/order';
import { ShipmentShape } from 'types/shipment';
import determineShipmentInfo from 'utils/shipmentInfo';

export class SelectShipmentType extends Component {
  constructor(props) {
    super(props);
    this.state = {
      showStorageInfoModal: false,
      showMoveInfoModal: false,
      errorMessage: null,
    };
  }

  componentDidMount() {
    const { loadMTOShipments, move } = this.props;
    loadMTOShipments(move.id);
  }

  setMoveType = (e) => {
    this.setState({ moveType: e.target.value });
  };

  toggleStorageModal = () => {
    this.setState((state) => ({
      showStorageInfoModal: !state.showStorageInfoModal,
    }));
  };

  toggleMoveInfoModal = () => {
    this.setState((state) => ({
      showMoveInfoModal: !state.showMoveInfoModal,
    }));
  };

  handleSubmit = () => {
    const { push, move, updateMove } = this.props;
    const { moveType } = this.state;

    const createShipmentPath = generatePath(customerRoutes.SHIPMENT_CREATE_PATH, { moveId: move.id });
    return patchMove({
      id: move.id,
      selected_move_type: moveType,
    })
      .then((response) => {
        // Update Redux with new data
        updateMove(response);
        push(`${createShipmentPath}?type=${moveType}`);
      })
      .catch((e) => {
        const { response } = e;
        const errorMessage = getResponseError(response, 'failed to update move due to server error');
        this.setState({
          errorMessage,
        });
      });
  };

  render() {
    const { push, move, mtoShipments } = this.props;
    const { moveType, showStorageInfoModal, showMoveInfoModal, errorMessage } = this.state;

    const shipmentInfo = determineShipmentInfo(move, mtoShipments);

    const ppmCardText = shipmentInfo.isPPMSelectable
      ? 'You pack and move your things, or make other arrangements, The government pays you for the weight you move.  This is a Personally Procured Move (PPM), sometimes called a DITY.'
      : 'You’ve already requested a PPM shipment. If you have more things to move yourself but that you can’t add to that shipment, contact the PPPO at your origin duty location.';

    const hhgCardText = shipmentInfo.isHHGSelectable
      ? 'All of your things are packed and moved by professionals, paid for by the government. This is a Household Goods move (HHG).'
      : 'Talk with your movers directly if you want to add or change shipments.';

    const ntsCardText = shipmentInfo.isNTSSelectable
      ? `Movers pack and ship things to a storage facility, where they stay until a future move. This is an NTS (non-temporary storage) shipment.`
      : 'You’ve already requested a long-term storage shipment for this move. Talk to your movers to change or add to your request.';

    const ntsrCardText = shipmentInfo.isNTSRSelectable
      ? 'Movers pick up things you put into NTS during an earlier move and ship them to your new destination. This is an NTS-release (non-temporary storage release) shipment.'
      : 'You’ve already asked to have things taken out of storage for this move. Talk to your movers to change or add to your request.';

    const selectableCardDefaultProps = {
      onChange: (e) => this.setMoveType(e),
      name: 'moveType',
    };

    const handleBack = () => {
      const backPath = shipmentInfo.hasShipment
        ? generalRoutes.HOME_PATH
        : generatePath(customerRoutes.SHIPMENT_MOVING_INFO_PATH, { moveId: move.id });

      push(backPath);
    };

    return (
      <>
        <GridContainer>
          <ScrollToTop otherDep={errorMessage} />

          <Grid row>
            <Grid col desktop={{ col: 8, offset: 2 }}>
              {errorMessage && (
                <Alert type="error" headingLevel="h4" heading="An error occurred">
                  {errorMessage}
                </Alert>
              )}

              <h6 data-testid="number-eyebrow" className="sm-heading margin-top-205 margin-bottom-0">
                Shipment {shipmentInfo.shipmentNumber}
              </h6>

              <h1 className={`${styles.selectTypeHeader} ${styles.header}`} data-testid="select-move-type-header">
                How should this shipment move?
              </h1>

              <p>
                You can move everything in one shipment, or you can split your belongings into multiple shipments that
                are moved in different ways.
              </p>
              <p>After you set up this shipment, you can add another shipment if you have more things to move.</p>

              <SelectableCard
                {...selectableCardDefaultProps}
                label="Movers pack and ship it, paid by the government (HHG)"
                value={SHIPMENT_OPTIONS.HHG}
                id={SHIPMENT_OPTIONS.HHG}
                cardText={hhgCardText}
                checked={moveType === SHIPMENT_OPTIONS.HHG}
                disabled={!shipmentInfo.isHHGSelectable}
                onHelpClick={this.toggleMoveInfoModal}
              />
              <SelectableCard
                {...selectableCardDefaultProps}
                label="Move it yourself and get paid for it (PPM)"
                value={SHIPMENT_OPTIONS.PPM}
                id={SHIPMENT_OPTIONS.PPM}
                cardText={ppmCardText}
                checked={moveType === SHIPMENT_OPTIONS.PPM}
                disabled={!shipmentInfo.isPPMSelectable}
                onHelpClick={this.toggleMoveInfoModal}
              />

              <h3 className={styles.longTermStorageHeader} data-testid="long-term-storage-heading">
                Long-term storage
              </h3>

              {!shipmentInfo.isNTSSelectable && !shipmentInfo.isNTSRSelectable ? (
                <p className={styles.pSmall}>
                  Talk to your movers about long-term storage if you need to add it to this move or change a request you
                  made earlier.
                </p>
              ) : (
                <>
                  <p>Your orders might not authorize long-term storage &mdash; your counselor can verify.</p>
                  <SelectableCard
                    {...selectableCardDefaultProps}
                    label="It's going into storage for months or years (NTS)"
                    value={SHIPMENT_OPTIONS.NTS}
                    id={SHIPMENT_OPTIONS.NTS}
                    cardText={ntsCardText}
                    checked={moveType === SHIPMENT_OPTIONS.NTS && shipmentInfo.isNTSSelectable}
                    disabled={!shipmentInfo.isNTSSelectable}
                    onHelpClick={this.toggleStorageModal}
                  />
                  <SelectableCard
                    {...selectableCardDefaultProps}
                    label="It was stored during a previous move (NTS-release)"
                    value={SHIPMENT_OPTIONS.NTSR}
                    id={SHIPMENT_OPTIONS.NTSR}
                    cardText={ntsrCardText}
                    checked={moveType === SHIPMENT_OPTIONS.NTSR && shipmentInfo.isNTSRSelectable}
                    disabled={!shipmentInfo.isNTSRSelectable}
                    onHelpClick={this.toggleStorageModal}
                  />
                </>
              )}

              {!shipmentInfo.hasShipment && (
                <p data-testid="helper-footer" className={styles.footer}>
                  <small>
                    It’s OK if you’re not sure about your choices. Your move counselor will go over all your options and
                    can help make changes if necessary.
                  </small>
                </p>
              )}

              <div className={formStyles.formActions}>
                <WizardNavigation
                  disableNext={moveType === undefined || moveType === ''}
                  onBackClick={handleBack}
                  onNextClick={this.handleSubmit}
                />
              </div>
            </Grid>
          </Grid>
        </GridContainer>
        <ConnectedMoveInfoModal isOpen={showMoveInfoModal} closeModal={this.toggleMoveInfoModal} />
        <ConnectedStorageInfoModal isOpen={showStorageInfoModal} closeModal={this.toggleStorageModal} />
      </>
    );
  }
}

SelectShipmentType.propTypes = {
  push: func.isRequired,
  updateMove: func.isRequired,
  loadMTOShipments: func.isRequired,
  move: MoveTaskOrderShape.isRequired,
  mtoShipments: arrayOf(ShipmentShape).isRequired,
};

const mapStateToProps = (state) => {
  const move = selectCurrentMove(state) || {};
  const mtoShipments = selectMTOShipmentsForCurrentMove(state);

  return {
    move,
    mtoShipments,
  };
};

const mapDispatchToProps = {
  updateMove: updateMoveAction,
  loadMTOShipments: loadMTOShipmentsAction,
};

export default connect(mapStateToProps, mapDispatchToProps)(SelectShipmentType);
export { mapStateToProps as _mapStateToProps };
