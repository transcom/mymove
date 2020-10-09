// NOTE: mock setup has to happen before import of file with dep that needs mocking

const mockStore = {};
let mockMove = {
  id: 'mockId',
  personally_procured_moves: ['mockPpm'],
  status: 'DRAFT',
};
let mockMtoShipments = ['mockHhgShipment'];

jest.mock('shared/Entities/modules/moves', () => {
  return {
    __esModule: true,
    selectActiveOrLatestMove: () => mockMove,
  };
});

jest.mock('shared/Entities/modules/mtoShipments', () => {
  return {
    __esModule: true,
    selectMTOShipmentsByMoveId: () => mockMtoShipments,
  };
});

// importing method to test now
import { _mapStateToProps } from 'pages/MyMove/SelectMoveType'; // eslint-disable-line

describe("SelectMoveType's mapStateToProps", () => {
  beforeEach(() => {
    // reset mock returns from internal dependencies
    mockMove = {
      id: 'mockId',
      personally_procured_moves: ['mockPpm'],
      status: 'DRAFT',
    };
    mockMtoShipments = ['mockHhgShipment'];
  });

  it('should set isHhgSelectable to false if move is already submitted', () => {
    // confirm initial state
    let actualProps = _mapStateToProps(mockStore);
    expect(actualProps.isHhgSelectable).toBe(true);

    // now validate test assertion
    mockMove.status = 'SUBMITTED';
    actualProps = _mapStateToProps(mockStore);
    expect(actualProps.isHhgSelectable).toBe(false);
  });

  it('should set isPpmSelectable to true if move does not have a PPM, even if the move is already submitted', () => {
    mockMove.personally_procured_moves = [];
    const actualProps = _mapStateToProps(mockStore);
    expect(actualProps.isPpmSelectable).toBe(true);
  });

  it('should set isPpmSelectable to false if move already has a PPM', () => {
    const actualProps = _mapStateToProps(mockStore);
    expect(actualProps.isPpmSelectable).toBe(false);
  });

  it('should return the correct new shipment number if a move has only PPM', () => {
    // clear mtoShipments from mock dep returns
    mockMtoShipments = [];

    const actualProps = _mapStateToProps(mockStore);
    expect(actualProps.shipmentNumber).toBe(2);
  });

  it('should return the correct new shipment number if a move has only HHG shipments', () => {
    // clear ppm from mock dep returns
    mockMove.personally_procured_moves = [];
    let actualProps = _mapStateToProps(mockStore);
    expect(actualProps.shipmentNumber).toBe(2);

    // let's add another mtoShipment to confirm it still gives the right number
    mockMtoShipments.push('anotherMockShipment');
    actualProps = _mapStateToProps(mockStore);
    expect(actualProps.shipmentNumber).toBe(3);
  });

  it('should return the correct new shipment number if a move has both PPM and HHG shipments', () => {
    let actualProps = _mapStateToProps(mockStore);
    expect(actualProps.shipmentNumber).toBe(3);

    // let's add another mtoShipment to confirm it still gives the right number
    mockMtoShipments.push('anotherMockShipment');
    actualProps = _mapStateToProps(mockStore);
    expect(actualProps.shipmentNumber).toBe(4);
  });
});
