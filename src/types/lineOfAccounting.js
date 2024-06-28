/* eslint-disable import/prefer-default-export */
import PropTypes from 'prop-types';

export const LineOfAccountingShape = PropTypes.shape({
  ID: PropTypes.string,
  LoaActvtyID: PropTypes.string,
  LoaAgncAcntngCd: PropTypes.string,
  LoaAgncDsbrCd: PropTypes.string,
  LoaAlltSnID: PropTypes.string,
  LoaBafID: PropTypes.string,
  LoaBdgtAcntClsNm: PropTypes.string,
  LoaBetCd: PropTypes.string,
  LoaBgFyTx: PropTypes.number,
  LoaBgnDt: PropTypes.string,
  LoaBgtLnItmID: PropTypes.string,
  LoaBgtRstrCd: PropTypes.string,
  LoaBgtSubActCd: PropTypes.string,
  LoaClsRefID: PropTypes.string,
  LoaCstCd: PropTypes.string,
  LoaCstCntrID: PropTypes.string,
  LoaCustNm: PropTypes.string,
  LoaDfAgncyAlctnRcpntID: PropTypes.string,
  LoaDocID: PropTypes.string,
  LoaDptID: PropTypes.string,
  LoaDscTx: PropTypes.string,
  LoaDtlRmbsmtSrcID: PropTypes.string,
  LoaEndDt: PropTypes.string,
  LoaEndFyTx: PropTypes.number,
  LoaFmsTrnsactnID: PropTypes.string,
  LoaFnclArID: PropTypes.string,
  LoaFnctPrsNm: PropTypes.string,
  LoaFndCntrID: PropTypes.string,
  LoaFndTyFgCd: PropTypes.string,
  LoaHistStatCd: PropTypes.string,
  LoaHsGdsCd: PropTypes.string,
  LoaInstlAcntgActID: PropTypes.string,
  LoaJbOrdNm: PropTypes.string,
  LoaLclInstlID: PropTypes.string,
  LoaMajClmNm: PropTypes.string,
  LoaMajRmbsmtSrcID: PropTypes.string,
  LoaObjClsID: PropTypes.string,
  LoaOpAgncyID: PropTypes.string,
  LoaPgmElmntID: PropTypes.string,
  LoaPrjID: PropTypes.string,
  LoaSbaltmtRcpntID: PropTypes.string,
  LoaScrtyCoopCustCd: PropTypes.string,
  LoaScrtyCoopDsgntrCd: PropTypes.string,
  LoaScrtyCoopImplAgncCd: PropTypes.string,
  LoaScrtyCoopLnItmID: PropTypes.string,
  LoaSpclIntrID: PropTypes.string,
  LoaSrvSrcID: PropTypes.string,
  LoaStatCd: PropTypes.string,
  LoaSubAcntID: PropTypes.string,
  LoaSysID: PropTypes.string,
  LoaTnsfrDptNm: PropTypes.string,
  LoaTrnsnID: PropTypes.string,
  LoaTrsySfxTx: PropTypes.string,
  LoaTskBdgtSblnTx: PropTypes.string,
  LoaUic: PropTypes.string,
  LoaWkCntrRcpntNm: PropTypes.string,
  LoaWrkOrdID: PropTypes.string,
  OrgGrpDfasCd: PropTypes.string,
  CreatedAt: PropTypes.string,
  UpdatedAt: PropTypes.string,
  ValidHhgProgramCodeForLoa: PropTypes.bool,
  ValidLoaForTac: PropTypes.bool,
});

// This is the order of DFAS elements in which to be concatenated into
// a "long line of accounting"
export const LineOfAccountingDfasElementOrder = [
  'loaDptID', // A1
  'loaTnsfrDptNm', // A2
  'loaEndFyTx', // A3
  'loaBafID', // A4
  'loaTrsySfxTx', // A5
  'loaMajClmNm', // A6
  'loaOpAgncyID', // B1
  'loaAlltSnID', // B2
  'loaUic', // B3
  'loaPgmElmntID', // C1
  'loaTskBdgtSblnTx', // C2
  'loaDfAgncyAlctnRcpntID', // D1
  'loaJbOrdNm', // D4
  'loaSbaltmtRcpntID', // D6
  'loaWkCntrRcpntNm', // D7
  'loaMajRmbsmtSrcID', // E1
  'loaDtlRmbsmtSrcID', // E2
  'loaCustNm', // E3
  'loaObjClsID', // F1
  'loaSrvSrcID', // F3
  'loaSpclIntrID', // G2
  'loaBdgtAcntClsNm', // I1
  'loaDocID', // J1
  'loaClsRefID', // K6
  'loaInstlAcntgActID', // L1
  'loaLclInstlID', // M1
  'loaTrnsnID', // N1
  'loaFmsTrnsactnID', // P5
];
