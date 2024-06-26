import numpy as np
import cv2
import random

# My library:
from opencv_transform.annotation import BodyPart


###
#
# maskdet_to_maskfin
#
# steps:
# 1. Extract annotation
# 1.a: Filter by color
# 1.b: Find ellipses
# 1.c: Filter out ellipses by max size, and max total numbers
# 1.d: Detect Problems
# 1.e: Resolve the problems, or discard the transformation
# 2. With the body list, draw maskfin, using maskref
#
###

# create_maskfin ==============================================================================
# return: (<Boolean> True/False), depending on the transformation process
def create_maskfin(maskref, maskdet):
    # Create a total green image, in which draw details ellipses
    details = np.zeros((512, 512, 3), np.uint8)
    details[:, :, :] = (0, 255, 0)  # (B, G, R)

    # Extract body part features:
    body_part_list = extract_annotations(maskdet)

    # Check if the list is not empty:
    if body_part_list:

        # Draw body part in details image:
        for obj in body_part_list:

            if obj.w < obj.h:
                a_max = int(obj.h / 2)  # asse maggiore
                a_min = int(obj.w / 2)  # asse minore
                angle = 0  # angle
            else:
                a_max = int(obj.w / 2)
                a_min = int(obj.h / 2)
                angle = 90

            x = int(obj.x)
            y = int(obj.y)

            # Draw ellipse
            if obj.name == "tit":
                cv2.ellipse(details, (x, y), (a_max, a_min), angle, 0, 360, (0, 205, 0), -1)  # (0,0,0,50)
            elif obj.name == "aur":
                cv2.ellipse(details, (x, y), (a_max, a_min), angle, 0, 360, (0, 0, 255), -1)  # red
            elif obj.name == "nip":
                cv2.ellipse(details, (x, y), (a_max, a_min), angle, 0, 360, (255, 255, 255), -1)  # white
            elif obj.name == "belly":
                cv2.ellipse(details, (x, y), (a_max, a_min), angle, 0, 360, (255, 0, 255), -1)  # purple
            elif obj.name == "vag":
                cv2.ellipse(details, (x, y), (a_max, a_min), angle, 0, 360, (255, 0, 0), -1)  # blue
            elif obj.name == "hair":
                x_min = x - int(obj.w / 2)
                y_min = y - int(obj.h / 2)
                x_max = x + int(obj.w / 2)
                y_max = y + int(obj.h / 2)
                cv2.rectangle(details, (x_min, y_min), (x_max, y_max), (100, 100, 100), -1)

        # Define the green color filter
        f1 = np.asarray([0, 250, 0])  # green color filter
        f2 = np.asarray([10, 255, 10])

        # From maskref, extrapolate only the green mask
        green_mask = cv2.bitwise_not(cv2.inRange(maskref, f1, f2))  # green is 0

        # Create an inverted mask
        green_mask_inv = cv2.bitwise_not(green_mask)

        # Cut maskref and detail image, using the green_mask & green_mask_inv
        res1 = cv2.bitwise_and(maskref, maskref, mask=green_mask)
        res2 = cv2.bitwise_and(details, details, mask=green_mask_inv)

        # Compone:
        maskfin = cv2.add(res1, res2)
        return maskfin


# extractAnnotations ==============================================================================
# input parameter:
# 	(<string> maskdet_img): relative path of the single maskdet image (es: testimg1/maskdet/1.png)
# return: (<BodyPart []> body_part_list) - for failure/error, return an empty list []
def extract_annotations(maskdet):
    # Load the image
    # image = cv2.imread(maskdet_img)

    # Find body part
    tits_list = find_body_part(maskdet, "tit")
    aur_list = find_body_part(maskdet, "aur")
    vag_list = find_body_part(maskdet, "vag")
    belly_list = find_body_part(maskdet, "belly")

    # Filter out parts basing on dimension (area and aspect ratio):
    aur_list = filter_dim_parts(aur_list, 100, 1000, 0.5, 3)
    tits_list = filter_dim_parts(tits_list, 1000, 60000, 0.2, 3)
    vag_list = filter_dim_parts(vag_list, 10, 1000, 0.2, 3)
    belly_list = filter_dim_parts(belly_list, 10, 1000, 0.2, 3)

    # Filter couple (if parts are > 2, choose only 2)
    aur_list = filter_couple(aur_list)
    tits_list = filter_couple(tits_list)

    # Detect a missing problem:
    missing_problem = detect_tit_aur_missing_problem(tits_list, aur_list)  # return a Number (code of the problem)

    # Check if problem is SOLVABLE:
    if missing_problem in [3, 6, 7, 8]:
        resolve_tit_aur_missing_problems(tits_list, aur_list, missing_problem)

    # Infer the nips:
    nip_list = infer_nip(aur_list)

    # Infer the hair:
    hair_list = infer_hair(vag_list)

    # Return a combined list:
    return tits_list + aur_list + nip_list + vag_list + hair_list + belly_list


# findBodyPart ==============================================================================
# input parameters:
# 	(<RGB>image, <string>part_name)
# return (<BodyPart[]>list)
def find_body_part(image, part_name):
    body_part_list = []  # empty BodyPart list
    color_mask = None

    # Get the correct color filter:
    if part_name == "tit":
        # Use combined color filter
        f1 = np.asarray([0, 0, 0])  # tit color filter
        f2 = np.asarray([10, 10, 10])
        f3 = np.asarray([0, 0, 250])  # aur color filter
        f4 = np.asarray([0, 0, 255])
        color_mask1 = cv2.inRange(image, f1, f2)
        color_mask2 = cv2.inRange(image, f3, f4)
        color_mask = cv2.bitwise_or(color_mask1, color_mask2)  # combine

    elif part_name == "aur":
        f1 = np.asarray([0, 0, 250])  # aur color filter
        f2 = np.asarray([0, 0, 255])
        color_mask = cv2.inRange(image, f1, f2)

    elif part_name == "vag":
        f1 = np.asarray([250, 0, 0])  # vag filter
        f2 = np.asarray([255, 0, 0])
        color_mask = cv2.inRange(image, f1, f2)

    elif part_name == "belly":
        f1 = np.asarray([250, 0, 250])  # belly filter
        f2 = np.asarray([255, 0, 255])
        color_mask = cv2.inRange(image, f1, f2)

    # find contours:
    contours, hierarchy = cv2.findContours(color_mask, cv2.RETR_TREE, cv2.CHAIN_APPROX_SIMPLE)

    # for every contour:
    for cnt in contours:

        if len(cnt) > 5:  # at least 5 points to fit ellipse

            # (x, y), (MA, ma), angle = cv2.fitEllipse(cnt)
            ellipse = cv2.fitEllipse(cnt)

            # Fit Result:
            x = ellipse[0][0]  # center x
            y = ellipse[0][1]  # center y
            angle = ellipse[2]  # angle
            a_min = ellipse[1][0]  # asse minore
            a_max = ellipse[1][1]  # asse maggiore

            # Detect direction:
            if angle == 0:
                h = a_max
                w = a_min
            else:
                h = a_min
                w = a_max

            # Normalize the belly size:
            if part_name == "belly":
                if w < 15:
                    w *= 2
                if h < 15:
                    h *= 2

            # Normalize the vag size:
            if part_name == "vag":
                if w < 15:
                    w *= 2
                if h < 15:
                    h *= 2

            # Calculate Bounding Box:
            x_min = int(x - (w / 2))
            x_max = int(x + (w / 2))
            y_min = int(y - (h / 2))
            y_max = int(y + (h / 2))

            body_part_list.append(BodyPart(part_name, x_min, y_min, x_max, y_max, x, y, w, h))

    return body_part_list


# filterDimParts ==============================================================================
# input parameters:
# 	(<BodyPart[]>list, <num> minimum area of part,  <num> max area, <num> min aspect ratio, <num> max aspect ratio)
def filter_dim_parts(bp_list, min_area, max_area, min_ar, max_ar):
    b_filt = []

    for obj in bp_list:

        a = obj.w * obj.h  # Object AREA

        if (a > min_area) and (a < max_area):

            ar = obj.w / obj.h  # Object ASPECT RATIO

            if (ar > min_ar) and (ar < max_ar):
                b_filt.append(obj)

    return b_filt


# filterCouple ==============================================================================
# input parameters:
# 	(<BodyPart[]>list)
def filter_couple(bp_list):
    # Remove exceed parts
    if len(bp_list) > 2:

        # trovare coppia (a,b) che minimizza bp_list[a].y-bp_list[b].y
        min_a = 0
        min_b = 1
        min_diff = abs(bp_list[min_a].y - bp_list[min_b].y)

        for a in range(0, len(bp_list)):
            for b in range(0, len(bp_list)):
                # TODO: avoid repetition (1,0) (0,1)
                if a != b:
                    diff = abs(bp_list[a].y - bp_list[b].y)
                    if diff < min_diff:
                        min_diff = diff
                        min_a = a
                        min_b = b
        b_filt = [bp_list[min_a], bp_list[min_b]]

        return b_filt
    else:
        # No change
        return bp_list


# detectTitAurMissingProblem ==============================================================================
# input parameters:
# 	(<BodyPart[]> tits list, <BodyPart[]> aur list)
# return (<num> problem code)
#   TIT  |  AUR  |  code |  SOLVE?  |
#    0   |   0   |   1   |    NO    |
#    0   |   1   |   2   |    NO    |
#    0   |   2   |   3   |    YES   |
#    1   |   0   |   4   |    NO    |
#    1   |   1   |   5   |    NO    |
#    1   |   2   |   6   |    YES   |
#    2   |   0   |   7   |    YES   |
#    2   |   1   |   8   |    YES   |
def detect_tit_aur_missing_problem(tits_list, aur_list):
    t_len = len(tits_list)
    a_len = len(aur_list)

    if t_len == 0:
        if a_len == 0:
            return 1
        elif a_len == 1:
            return 2
        elif a_len == 2:
            return 3
        else:
            return -1
    elif t_len == 1:
        if a_len == 0:
            return 4
        elif a_len == 1:
            return 5
        elif a_len == 2:
            return 6
        else:
            return -1
    elif t_len == 2:
        if a_len == 0:
            return 7
        elif a_len == 1:
            return 8
        else:
            return -1
    else:
        return -1


# resolveTitAurMissingProblems ==============================================================================
# input parameters:
# 	(<BodyPart[]> tits list, <BodyPart[]> aur list, problem code)
# return None
def resolve_tit_aur_missing_problems(tits_list, aur_list, problem_code):
    if problem_code == 3:

        random_tit_factor = random.randint(2, 5)  # TO_TEST

        # Add the first tit:
        new_w = aur_list[0].w * random_tit_factor  # TO_TEST
        new_x = aur_list[0].x
        new_y = aur_list[0].y

        x_min = int(new_x - (new_w / 2))
        x_max = int(new_x + (new_w / 2))
        y_min = int(new_y - (new_w / 2))
        y_max = int(new_y + (new_w / 2))

        tits_list.append(BodyPart("tit", x_min, y_min, x_max, y_max, new_x, new_y, new_w, new_w))

        # Add the second tit:
        new_w = aur_list[1].w * random_tit_factor  # TO_TEST
        new_x = aur_list[1].x
        new_y = aur_list[1].y

        x_min = int(new_x - (new_w / 2))
        x_max = int(new_x + (new_w / 2))
        y_min = int(new_y - (new_w / 2))
        y_max = int(new_y + (new_w / 2))

        tits_list.append(BodyPart("tit", x_min, y_min, x_max, y_max, new_x, new_y, new_w, new_w))

    elif problem_code == 6:

        # Find wich aur is full:
        d1 = abs(tits_list[0].x - aur_list[0].x)
        d2 = abs(tits_list[0].x - aur_list[1].x)

        if d1 > d2:
            # aur[0] is empty
            new_x = aur_list[0].x
            new_y = aur_list[0].y
        else:
            # aur[1] is empty
            new_x = aur_list[1].x
            new_y = aur_list[1].y

        # Calculate Bounding Box:
        x_min = int(new_x - (tits_list[0].w / 2))
        x_max = int(new_x + (tits_list[0].w / 2))
        y_min = int(new_y - (tits_list[0].w / 2))
        y_max = int(new_y + (tits_list[0].w / 2))

        tits_list.append(BodyPart("tit", x_min, y_min, x_max, y_max, new_x, new_y, tits_list[0].w, tits_list[0].w))

    elif problem_code == 7:

        # Add the first aur:
        new_w = tits_list[0].w * random.uniform(0.03, 0.1)  # TO_TEST
        new_x = tits_list[0].x
        new_y = tits_list[0].y

        x_min = int(new_x - (new_w / 2))
        x_max = int(new_x + (new_w / 2))
        y_min = int(new_y - (new_w / 2))
        y_max = int(new_y + (new_w / 2))

        aur_list.append(BodyPart("aur", x_min, y_min, x_max, y_max, new_x, new_y, new_w, new_w))

        # Add the second aur:
        new_w = tits_list[1].w * random.uniform(0.03, 0.1)  # TO_TEST
        new_x = tits_list[1].x
        new_y = tits_list[1].y

        x_min = int(new_x - (new_w / 2))
        x_max = int(new_x + (new_w / 2))
        y_min = int(new_y - (new_w / 2))
        y_max = int(new_y + (new_w / 2))

        aur_list.append(BodyPart("aur", x_min, y_min, x_max, y_max, new_x, new_y, new_w, new_w))

    elif problem_code == 8:

        # Find wich tit is full:
        d1 = abs(aur_list[0].x - tits_list[0].x)
        d2 = abs(aur_list[0].x - tits_list[1].x)

        if d1 > d2:
            # tit[0] is empty
            new_x = tits_list[0].x
            new_y = tits_list[0].y
        else:
            # tit[1] is empty
            new_x = tits_list[1].x
            new_y = tits_list[1].y

        # Calculate Bounding Box:
        x_min = int(new_x - (aur_list[0].w / 2))
        x_max = int(new_x + (aur_list[0].w / 2))
        y_min = int(new_y - (aur_list[0].w / 2))
        y_max = int(new_y + (aur_list[0].w / 2))
        aur_list.append(BodyPart("aur", x_min, y_min, x_max, y_max, new_x, new_y, aur_list[0].w, aur_list[0].w))


# detectTitAurPositionProblem ==============================================================================
# input parameters:
# 	(<BodyPart[]> tits list, <BodyPart[]> aur list)
# return (<Boolean> True/False)
def detect_tit_aur_position_problem(tits_list, aur_list):
    diff_tits_x = abs(tits_list[0].x - tits_list[1].x)
    if diff_tits_x < 40:
        print("diff_tits_x")
        # Tits too narrow (horizontally)
        return True

    diff_tits_y = abs(tits_list[0].y - tits_list[1].y)
    if diff_tits_y > 120:
        # Tits too distanced (vertically)
        print("diff_tits_y")
        return True

    diff_tits_w = abs(tits_list[0].w - tits_list[1].w)
    if (diff_tits_w < 0.1) or (diff_tits_w > 60):
        print("diff_tits_w")
        # Tits too equals, or too different (width)
        return True

    # Check if body position is too low (face not covered by watermark)
    if aur_list[0].y > 350:  # tits too low
        # Calculate the ratio between y and aurs distance
        rapp = aur_list[0].y / (abs(aur_list[0].x - aur_list[1].x))
        if rapp > 2.8:
            print("aurDown")
            return True

    return False


# infer_nip ==============================================================================
# input parameters:
# 	(<BodyPart[]> aur list)
# return (<BodyPart[]> nip list)
def infer_nip(aur_list):
    nip_list = []

    for aur in aur_list:
        # Nip rules:
        # - circle (w == h)
        # - min dim: 5
        # - bigger if aur is bigger
        nip_dim = int(5 + aur.w * random.uniform(0.03, 0.09))

        # center:
        x = aur.x
        y = aur.y

        # Calculate Bounding Box:
        x_min = int(x - (nip_dim / 2))
        x_max = int(x + (nip_dim / 2))
        y_min = int(y - (nip_dim / 2))
        y_max = int(y + (nip_dim / 2))

        nip_list.append(BodyPart("nip", x_min, y_min, x_max, y_max, x, y, nip_dim, nip_dim))

    return nip_list


# infer_hair (TO_TEST) ==============================================================================
# input parameters:
# 	(<BodyPart[]> vag list)
# return (<BodyPart[]> hair list)
def infer_hair(vag_list):
    hair_list = []

    # 70% of chance to add hair
    if random.uniform(0.0, 1.0) > 0.3:

        for vag in vag_list:
            # Hair rules:
            hair_w = vag.w * random.uniform(0.4, 1.5)
            hair_h = vag.h * random.uniform(0.4, 1.5)

            # center:
            x = vag.x
            y = vag.y - (hair_h / 2) - (vag.h / 2)

            # Calculate Bounding Box:
            x_min = int(x - (hair_w / 2))
            x_max = int(x + (hair_w / 2))
            y_min = int(y - (hair_h / 2))
            y_max = int(y + (hair_h / 2))

            hair_list.append(BodyPart("hair", x_min, y_min, x_max, y_max, x, y, hair_w, hair_h))

    return hair_list
