package notacore

type Input int

const (
	InputInvalid Input = iota

	KeySpace
	KeyApostrophe
	KeyComma
	KeyMinus
	KeyPeriod
	KeySlash

	Key0
	Key1
	Key2
	Key3
	Key4
	Key5
	Key6
	Key7
	Key8
	Key9

	KeySemicolon
	KeyEqual

	KeyA
	KeyB
	KeyC
	KeyD
	KeyE
	KeyF
	KeyG
	KeyH
	KeyI
	KeyJ
	KeyK
	KeyL
	KeyM
	KeyN
	KeyO
	KeyP
	KeyQ
	KeyR
	KeyS
	KeyT
	KeyU
	KeyV
	KeyW
	KeyX
	KeyY
	KeyZ

	KeyLeftBracket
	KeyBackslash
	KeyRightBracket
	KeyGraveAccent

	KeyEscape
	KeyEnter
	KeyTab
	KeyBackspace
	KeyInsert
	KeyDelete

	KeyRight
	KeyLeft
	KeyDown
	KeyUp

	KeyPageUp
	KeyPageDown
	KeyHome
	KeyEnd

	KeyCapsLock
	KeyScrollLock
	KeyNumLock
	KeyPrintScreen
	KeyPause

	KeyF1
	KeyF2
	KeyF3
	KeyF4
	KeyF5
	KeyF6
	KeyF7
	KeyF8
	KeyF9
	KeyF10
	KeyF11
	KeyF12
	KeyF13
	KeyF14
	KeyF15
	KeyF16
	KeyF17
	KeyF18
	KeyF19
	KeyF20
	KeyF21
	KeyF22
	KeyF23
	KeyF24
	KeyF25

	KeyKP0
	KeyKP1
	KeyKP2
	KeyKP3
	KeyKP4
	KeyKP5
	KeyKP6
	KeyKP7
	KeyKP8
	KeyKP9
	KeyKPDecimal
	KeyKPDivide
	KeyKPMultiply
	KeyKPSubtract
	KeyKPAdd
	KeyKPEnter
	KeyKPEqual

	KeyLeftShift
	KeyLeftControl
	KeyLeftAlt
	KeyLeftSuper
	KeyRightShift
	KeyRightControl
	KeyRightAlt
	KeyRightSuper
	KeyLeftCommand  //does not work yet
	KeyRightCommand // does not work yet
	KeyOptionLeft   //does not work yet
	KeyOptionRight  //does not work yet
	KeyFn           //does not work yet
	KeyMenu

	KeyMediaPlayPause
	KeyMediaStop
	KeyMediaNext
	KeyMediaPrev
	KeyVolumeUp
	KeyVolumeDown
	KeyBrightnessUp
	KeyBrightnessDown

	MouseLeft
	MouseRight
	MouseMiddle
	MouseButton4
	MouseButton5
	MouseButton6
	MouseButton7
	MouseButton8

	MouseX
	MouseY
	MouseScrollX
	MouseScrollY

	PadA
	PadB
	PadX
	PadY
	PadLB
	PadRB
	PadBack
	PadStart
	PadGuide
	PadLeftThumb
	PadRightThumb
	PadDpadUp
	PadDpadRight
	PadDpadDown
	PadDpadLeft

	PadAxisLeftX
	PadAxisLeftY
	PadAxisRightX
	PadAxisRightY
	PadAxisLeftTrigger
	PadAxisRightTrigger

	//TODO: what is below is not working yet, OS based API needed

	Touch1
	Touch2
	Touch3
	Touch4
	Touch5
	ForceTouch

	SwipeUp1
	SwipeDown1
	SwipeLeft1
	SwipeRight1
	SwipeUp2
	SwipeDown2
	SwipeLeft2
	SwipeRight2
	SwipeUp3
	SwipeDown3
	SwipeLeft3
	SwipeRight3
	SwipeUp4
	SwipeDown4
	SwipeLeft4
	SwipeRight4
	SwipeUp5
	SwipeDown5
	SwipeLeft5
	SwipeRight5

	PinchIn
	PinchOut
	RotateCW
	RotateCCW

	MobileTouchBegin
	MobileTouchMove
	MobileTouchEnd
	MobileTouchCancel
	MobileBack
	MobileHome
	MobileVolumeUp
	MobileVolumeDown
	AccelerometerX
	AccelerometerY
	AccelerometerZ
	GyroX
	GyroY
	GyroZ
	OrientationPitch
	OrientationYaw
	OrientationRoll
)
