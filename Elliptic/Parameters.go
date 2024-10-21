package Elliptic

import (
    aux "DALOS_Crypto/Auxilliary"
    "math/big"
    "strconv"
)

var (
    Zero = big.NewInt(0)
    One  = big.NewInt(1)
    Two  = big.NewInt(2)
    
    InfinityPoint = CoordExtended{Zero, One, One, Zero}
)

type Ellipse struct {
    Name string //Name
    //Prime Numbers
    P big.Int //Prime Number defining the Prime Field
    Q big.Int //Prime Number defining the Generator (Base-Point) Order
    T big.Int //Trace of the Curve
    R big.Int //Elliptic Curve Cofactor: R*Q = P + 1 - T
    
    //Coefficients (Equation Parameters)
    A big.Int // x^2 Coefficient (Twisted Edwards Curve)
    D big.Int // x^2 * y^2 Coefficient (Twisted Edwards Curve)
    
    //Curve safe scalar size in bits
    S uint32
    
    //Point Coordinates
    G CoordAffine
}

type PrimePowerTwo struct {
    Power      int
    RestString string
    Sign       bool
}

func MakePrime(PrimeNumber PrimePowerTwo) big.Int {
    var (
        Prime = new(big.Int)
        Rest  = new(big.Int)
    )
    
    Rest.SetString(PrimeNumber.RestString, 10)
    if PrimeNumber.Sign == true {
        Prime.SetBit(Zero, PrimeNumber.Power, 1).Add(Prime, Rest)
    } else {
        Prime.SetBit(Zero, PrimeNumber.Power, 1).Sub(Prime, Rest)
    }
    
    return *Prime
}

func ComputeCofactor(P, Q, T big.Int) big.Int {
    var h = new(big.Int)
    h.Add(&P, One).Sub(h, &T).Quo(h, &Q)
    return *h
}

func ComputeSafeScalar(Prime, Trace, Cofactor *big.Int) (uint64, string) {
    //@doc "Computes the Safe Scalar, given Prime Number, Trace and Cofactor of an Elliptic Curve
    //The Safe scalar, is the power of 2, in this case 2^1600 means this many private keys possible
    //Computing the Safe Scalar assumes, the Prime is a prime Number, the Cofactor is correct for the Elliptic Curve
    //and the Trace is also Correct. Computing the Safe Scalar also yields the Generator of the elliptic curve."
    var (
        Q         = new(big.Int)
        Remainder = new(big.Int)
        Qs        string
        Ss        string
        X         uint64
    )
    
    CofactorBase2 := Cofactor.Text(2)
    CofactorBase2Trimmed := aux.TrimFirstRune(CofactorBase2)
    CofactorBitSize := uint64(len(CofactorBase2Trimmed))
    v1 := InferiorTrace(Prime, Trace)
    Q.QuoRem(v1, Cofactor, Remainder)
    Power, Sign, Rest := Power2DistanceChecker(Q)
    
    if Sign == false {
        X = Power - (2 + CofactorBitSize)
        Ss = "-"
    } else if Sign == true {
        X = Power - (1 + CofactorBitSize)
        Ss = "+"
    }
    PowerString := strconv.FormatInt(int64(Power), 10)
    RestS := Rest.Text(10)
    Qs = "2^" + PowerString + Ss + RestS
    return X, Qs
}

func Power2DistanceChecker(Number *big.Int) (uint64, bool, *big.Int) {
    //@doc "Transforms a big.Int in 2^x +/- y representation,
    //Returning the Power(uint), Sign(bool), and the remainder, y(*big.Int)"
    var (
        BetweenInt  = new(big.Int)
        HalfBetween = new(big.Int)
        LowerPower  = new(big.Int)
        HigherPower = new(big.Int)
        
        Rest  = new(big.Int)
        Sign  bool
        Power uint64
    )
    NumberBase2 := Number.Text(2)
    Between := aux.TrimFirstRune(NumberBase2)
    BetweenInt.SetString(Between, 2)   //22
    LowerPower.Sub(Number, BetweenInt) //32
    HigherPower.Mul(LowerPower, Two)   //64
    HalfBetween.Quo(LowerPower, Two)   //16
    HigherPowerBin := HigherPower.Text(2)
    Cmp := BetweenInt.Cmp(HalfBetween)
    if Cmp == 1 {
        Rest.Sub(LowerPower, BetweenInt)
        Sign = false
        Power = uint64(len(HigherPowerBin)) - 1
    } else if Cmp == -1 {
        Rest = BetweenInt
        Sign = true
        Power = uint64(len(HigherPowerBin)) - 2
    } else {
        Rest = BetweenInt
        Sign = true
        Power = uint64(len(HigherPowerBin)) - 2
    }
    return Power, Sign, Rest
}

func InferiorTrace(Prime, Trace *big.Int) *big.Int {
    var output = new(big.Int)
    return output.Add(Prime, One).Sub(output, Trace)
}

func SuperiorTrace(Prime, Trace *big.Int) *big.Int {
    var output = new(big.Int)
    return output.Add(Prime, One).Add(output, Trace)
}

func E521Ellipse() Ellipse {
    var (
        e Ellipse
        P PrimePowerTwo
        Q PrimePowerTwo
    )
    //Name
    e.Name = "E521"
    
    //Prime Numbers
    //P=2^521 - 1
    P.Power = 521
    P.RestString = "1"
    P.Sign = false
    e.P = MakePrime(P)
    
    //Q=2^519 - 337554763258501705789107630418782636071904961214051226618635150085779108655765
    Q.Power = 519
    Q.RestString = "337554763258501705789107630418782636071904961214051226618635150085779108655765"
    Q.Sign = false
    e.Q = MakePrime(Q)
    
    //Trace and Cofactor
    e.T.SetString("1350219053034006823156430521675130544287619844856204906474540600343116434623060", 10)
    e.R = ComputeCofactor(e.P, e.Q, e.T)
    
    //Safe Scalar Size in bits = 1600
    e.S = 515
    
    //A and D Coefficients
    e.A.SetInt64(1)
    e.D.SetInt64(-376014)
    
    //Generator Coordinates
    e.G.AX.SetString("1571054894184995387535939749894317568645297350402905821437625181152304994381188529632591196067604100772673927915114267193389905003276673749012051148356041324", 10)
    e.G.AY.SetInt64(12)
    
    return e
}

func DalosEllipse() Ellipse {
    var (
        e Ellipse
        P PrimePowerTwo
        Q PrimePowerTwo
    )
    //Name
    e.Name = "TEC_S1600_Pr1605p2315_m26"
    
    //Prime Numbers
    //P=2^1605 + 2315
    P.Power = 1605
    P.RestString = "2315"
    P.Sign = true
    e.P = MakePrime(P)
    
    //2^1603+1258387060301909514024042379046449850251725029634697115619073843890705481440046740552204199635883885272944914904655483501916023678206167596650367826811846862157534952990004386839463386963494516862067933899764941962204635259228497801901380413
    Q.Power = 1603
    Q.RestString = "1258387060301909514024042379046449850251725029634697115619073843890705481440046740552204199635883885272944914904655483501916023678206167596650367826811846862157534952990004386839463386963494516862067933899764941962204635259228497801901380413"
    Q.Sign = true
    e.Q = MakePrime(Q)
    
    //Trace and Cofactor
    e.T.SetString("-5033548241207638056096169516185799401006900118538788462476295375562821925760186962208816798543535541091779659618621934007664094712824670386601471307247387448630139811960017547357853547853978067448271735599059767848818541036913991207605519336", 10)
    e.R = ComputeCofactor(e.P, e.Q, e.T)
    
    //Safe Scalar Size in bits = 1600
    e.S = 1600
    
    //A and D Coefficients
    e.A.SetInt64(1)
    e.D.SetInt64(-26)
    
    //Generator Coordinates
    e.G.AX = new(big.Int) // Allocate memory for AX
    e.G.AY = new(big.Int) // Allocate memory for AY
    e.G.AX.SetInt64(2)
    e.G.AY.SetString("479577721234741891316129314062096440203224800598561362604776518993348406897758651324205216647014453759416735508511915279509434960064559686580741767201752370055871770203009254182472722342456597752506165983884867351649283353392919401537107130232654743719219329990067668637876645065665284755295099198801899803461121192253205447281506198423683290960014859350933836516450524873032454015597501532988405894858561193893921904896724509904622632232182531698393484411082218273681226753590907472", 10)
    
    return e
}
