/* eslint-disable react/jsx-props-no-spreading */
import React, { Component } from 'react';
import { connect } from 'react-redux';
import { string, bool, func, arrayOf, shape } from 'prop-types';

import styles from './SelectMoveType.module.scss';

import { SHIPMENT_OPTIONS, MOVE_STATUSES } from 'shared/constants';
import { selectCurrentMove, selectMTOShipmentsForCurrentMove } from 'store/entities/selectors';
import { WizardPage } from 'shared/WizardPage';
import SelectableCard from 'components/Customer/SelectableCard';
import { loadMTOShipments as loadMTOShipmentsAction } from 'shared/Entities/modules/mtoShipments';
import { patchMove, getResponseError } from 'services/internalApi';
import { updateMove as updateMoveAction } from 'store/entities/actions';
import { MoveTaskOrderShape, MTOShipmentShape } from 'types/order';
import ConnectedStorageInfoModal from 'components/Customer/modals/StorageInfoModal/StorageInfoModal';
import ConnectedMoveInfoModal from 'components/Customer/modals/MoveInfoModal/MoveInfoModal';

export class SelectMoveType extends Component {
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
    const { match, updateMove } = this.props;
    const { moveType } = this.state;

    return patchMove({
      id: match.params.moveId,
      selected_move_type: moveType,
    })
      .then((response) => {
        // Update Redux with new data
        updateMove(response);
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
    const { pageKey, pageList, match, push, move, mtoShipments } = this.props;
    const { moveType, showStorageInfoModal, showMoveInfoModal, errorMessage } = this.state;
    const hasNTS = mtoShipments.some((shipment) => shipment.shipmentType === SHIPMENT_OPTIONS.NTS);
    const hasNTSR = mtoShipments.some((shipment) => shipment.shipmentType === SHIPMENT_OPTIONS.NTSR);
    const isMoveDraft = move.status === MOVE_STATUSES.DRAFT;
    const hasPpm = !!move.personally_procured_moves?.length;
    const isPpmSelectable = !hasPpm;
    const isHhgSelectable = isMoveDraft;
    const isNtsSelectable = isMoveDraft && !hasNTS;
    const isNtsrSelectable = isMoveDraft && !hasNTSR;
    const ppmCount = hasPpm ? 1 : 0;
    const mtosCount = mtoShipments?.length || 0;
    const shipmentNumber = 1 + ppmCount + mtosCount;
    const hasShipment = ppmCount + mtosCount > 0;
    const canMoveNext = moveType ? moveType !== '' : false;
    const ppmCardText =
      'You pack and move your things, or make other arrangements, The government pays you for the weight you move.  This is a a Personally Procured Move (PPM), sometimes called a DITY.';
    const hhgCardText =
      'Your things are packed and moved by professionals, paid for by the government. This is a Household Goods move (HHG).';
    const ntsCardText = `Movers pack and ship things to a storage facility, where they stay until a future move. Your orders might not authorize long-term storage — your counselor can verify. This is an NTS (non-temporary storage) shipment.`;
    const ntsrCardText =
      'Movers pick up things you put into NTS during an earlier move and ship them to your new destination. This is an NTS-R (non-temporary storage release) shipment.';
    const ntsDisabledText =
      'You’ve already requested a long-term storage shipment for this move. Talk to your movers to change or add to your request.';
    const ntsrDisabledText =
      'You’ve already asked to have things taken out of storage for this move. Talk to your movers to change or add to your request.';
    const hhgCardTextPostSubmit = 'Talk with your movers directly if you want to add or change shipments.';
    const ppmCardTextAlreadyChosen = `You’ve already requested a PPM shipment. If you have more things to move yourself but that you can’t add to that shipment, contact the PPPO at your origin duty station.`;
    const noLongTermStorageCardsText =
      'Talk to your movers about long-term storage if you need to add it to this move or change a request you made earlier.';
    const selectableCardDefaultProps = {
      onChange: (e) => this.setMoveType(e),
      name: 'moveType',
    };
    const ppmEnabledCard = (
      <SelectableCard
        {...selectableCardDefaultProps}
        label="Do it yourself"
        value={SHIPMENT_OPTIONS.PPM}
        id={SHIPMENT_OPTIONS.PPM}
        cardText={ppmCardText}
        checked={moveType === SHIPMENT_OPTIONS.PPM}
        disabled={false}
        onHelpClick={this.toggleMoveInfoModal}
      />
    );
    const ppmDisabledCard = (
      <SelectableCard
        {...selectableCardDefaultProps}
        label="Do it yourself (already chosen)"
        value={SHIPMENT_OPTIONS.PPM}
        id={SHIPMENT_OPTIONS.PPM}
        cardText={ppmCardTextAlreadyChosen}
        checked={false}
        disabled={!isPpmSelectable}
        onHelpClick={this.toggleMoveInfoModal}
      />
    );
    const hhgEnabledCard = (
      <SelectableCard
        {...selectableCardDefaultProps}
        label="Professional movers"
        value={SHIPMENT_OPTIONS.HHG}
        id={SHIPMENT_OPTIONS.HHG}
        cardText={hhgCardText}
        checked={moveType === SHIPMENT_OPTIONS.HHG}
        disabled={false}
        onHelpClick={this.toggleMoveInfoModal}
      />
    );
    const hhgDisabledCard = (
      <SelectableCard
        {...selectableCardDefaultProps}
        label="Professional movers"
        value={SHIPMENT_OPTIONS.HHG}
        id={SHIPMENT_OPTIONS.HHG}
        cardText={hhgCardTextPostSubmit}
        checked={false}
        disabled={!isHhgSelectable}
        onHelpClick={this.toggleMoveInfoModal}
      />
    );
    const footerText = (
      <div>
        {!hasShipment && (
          <div data-testid="helper-footer" className={`${styles.footer} grid-col-12`}>
            It’s OK if you’re not sure about your choices. Your move counselor will go over all your options and can
            help make changes if necessary.
          </div>
        )}
      </div>
    );
    return (
      <>
        <WizardPage
          pageKey={pageKey}
          match={match}
          pageList={pageList}
          dirty
          handleSubmit={this.handleSubmit}
          push={push}
          footerText={footerText}
          canMoveNext={canMoveNext}
          error={errorMessage}
        >
          <h6 data-testid="number-eyebrow" className="sm-heading">
            Shipment {shipmentNumber}
          </h6>
          <h1 className={`${styles.selectTypeHeader} ${styles.header}`} data-testid="select-move-type-header">
            {shipmentNumber > 1
              ? 'How do you want this group of things moved?'
              : 'How do you want to move your belongings?'}
          </h1>
          <h2>Choose 1 shipment at a time.</h2>
          <p>You can add more later</p>
          {isPpmSelectable ? ppmEnabledCard : ppmDisabledCard}
          {isHhgSelectable ? hhgEnabledCard : hhgDisabledCard}
          <h3 data-testid="long-term-storage-heading">Long-term storage</h3>
          {!isNtsSelectable && !isNtsrSelectable ? (
            <p className={styles.pSmall}>{noLongTermStorageCardsText}</p>
          ) : (
            <>
              <p>These shipments do count against your weight allowance for this move.</p>
              <SelectableCard
                {...selectableCardDefaultProps}
                label="Put things into long-term storage"
                value={SHIPMENT_OPTIONS.NTS}
                id={SHIPMENT_OPTIONS.NTS}
                cardText={isNtsSelectable ? ntsCardText : ntsDisabledText}
                checked={moveType === SHIPMENT_OPTIONS.NTS && isNtsSelectable}
                disabled={!isNtsSelectable}
                onHelpClick={this.toggleStorageModal}
              />
              <SelectableCard
                {...selectableCardDefaultProps}
                label="Get things out of long-term storage"
                value={SHIPMENT_OPTIONS.NTSR}
                id={SHIPMENT_OPTIONS.NTSR}
                cardText={isNtsrSelectable ? ntsrCardText : ntsrDisabledText}
                checked={moveType === SHIPMENT_OPTIONS.NTSR && isNtsrSelectable}
                disabled={!isNtsrSelectable}
                onHelpClick={this.toggleStorageModal}
              />
            </>
          )}
        </WizardPage>
        <ConnectedMoveInfoModal isOpen={showMoveInfoModal} closeModal={this.toggleMoveInfoModal} />
        <ConnectedStorageInfoModal isOpen={showStorageInfoModal} closeModal={this.toggleStorageModal} />
      </>
    );
  }
}

SelectMoveType.propTypes = {
  pageKey: string.isRequired,
  pageList: arrayOf(string).isRequired,
  match: shape({
    isExact: bool.isRequired,
    params: shape({
      moveId: string.isRequired,
    }),
    path: string.isRequired,
    url: string.isRequired,
  }).isRequired,
  push: func.isRequired,
  updateMove: func.isRequired,
  loadMTOShipments: func.isRequired,
  move: MoveTaskOrderShape.isRequired,
  mtoShipments: arrayOf(MTOShipmentShape).isRequired,
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

export default connect(mapStateToProps, mapDispatchToProps)(SelectMoveType);
export { mapStateToProps as _mapStateToProps };
